package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dadosjusbr/api/models"
	"github.com/dadosjusbr/storage"
	strModels "github.com/dadosjusbr/storage/models"
	"github.com/dadosjusbr/storage/repositories/database/mongo"
	"github.com/dadosjusbr/storage/repositories/database/postgres"
	"github.com/dadosjusbr/storage/repositories/fileStorage"
	"github.com/dadosjusbr/storage/repositories/interfaces"
	"github.com/gocarina/gocsv"
	"github.com/joho/godotenv"
	"github.com/newrelic/go-agent/v3/integrations/nrecho-v3"
	"github.com/newrelic/go-agent/v3/newrelic"

	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type config struct {
	Port   int    `envconfig:"PORT"`
	DBUrl  string `envconfig:"MONGODB_URI"`
	DBName string `envconfig:"MONGODB_NAME"`

	AwsS3Bucket string `envconfig:"AWS_S3_BUCKET" required:"true"`
	AwsRegion   string `envconfig:"AWS_REGION" required:"true"`

	// StorageDB config
	MongoURI    string `envconfig:"MONGODB_URI"`
	MongoDBName string `envconfig:"MONGODB_NAME"`
	MongoMICol  string `envconfig:"MONGODB_MICOL" required:"true"`
	MongoAgCol  string `envconfig:"MONGODB_AGCOL" required:"true"`
	MongoPkgCol string `envconfig:"MONGODB_PKGCOL" required:"true"`
	MongoRevCol string `envconfig:"MONGODB_REVCOL" required:"true"`

	// Omited fields
	EnvOmittedFields []string `envconfig:"ENV_OMITTED_FIELDS"`

	// Site env
	DadosJusURL    string `envconfig:"DADOSJUS_URL" required:"true"`
	PackageRepoURL string `envconfig:"PACKAGE_REPO_URL" required:"true"`

	// PostgresDB config
	PgUser     string `envconfig:"PG_USER"`
	PgPassword string `envconfig:"PG_PASSWORD"`
	PgDatabase string `envconfig:"PG_DATABASE"`
	PgHost     string `envconfig:"PG_HOST"`
	PgPort     string `envconfig:"PG_PORT"`

	// Query limit env
	SearchLimit   int `envconfig:"SEARCH_LIMIT"`
	DownloadLimit int `envconfig:"DOWNLOAD_LIMIT"`

	// Newrelic config
	NewRelicApp     string `envconfig:"NEWRELIC_APP_NAME"`
	NewRelicLicense string `envconfig:"NEWRELIC_LICENSE"`
}

var mgoClient *storage.Client
var pgClient *storage.Client
var loc *time.Location
var conf config
var postgresDb *PostgresDB
var sess *AwsSession

// newClient takes a config struct and creates a client to connect with DB and Cloud5
func newClient(db interfaces.IDatabaseRepository) (*storage.Client, error) {
	client, err := storage.NewClient(db, &fileStorage.S3Client{})
	if err != nil {
		return nil, fmt.Errorf("error creating storage.client: %q", err)
	}
	return client, nil
}

func newPostgresDB(c config) (*postgres.PostgresDB, error) {
	pgDb, err := postgres.NewPostgresDB(c.PgUser, c.PgPassword, c.PgDatabase, c.PgHost, c.PgPort)
	if err != nil {
		return nil, fmt.Errorf("error creating postgres DB client: %q", err)
	}
	return pgDb, nil
}

func newMongoDB(c config) (*mongo.DBClient, error) {
	if c.MongoMICol == "" || c.MongoAgCol == "" {
		return nil, fmt.Errorf("error creating storage client: db collections must not be empty. MI:\"%s\", AG:\"%s\", PKG:\"%s\"", c.MongoMICol, c.MongoAgCol, c.MongoPkgCol)
	}
	db, err := mongo.NewMongoDB(c.MongoURI, c.MongoDBName, c.MongoMICol, c.MongoAgCol, c.MongoPkgCol, c.MongoRevCol)
	if err != nil {
		return nil, fmt.Errorf("error creating mongo DB client: %q", err)
	}
	db.Collection(c.MongoMICol)
	return db, nil
}

func getTotalsOfAgencyYear(c echo.Context) error {
	year, err := strconv.Atoi(c.Param("ano"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d inválido", year))
	}
	aID := c.Param("orgao")
	agenciesMonthlyInfo, err := pgClient.Db.GetMonthlyInfo([]strModels.Agency{{ID: aID}}, year)
	if err != nil {
		log.Printf("[totals of agency year] error getting data for first screen(ano:%d, estado:%s):%q", year, aID, err)
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d ou orgao=%s inválidos", year, aID))
	}
	var monthTotalsOfYear []models.MonthTotals
	agency, err := pgClient.Db.GetAgency(aID)
	if err != nil {
		log.Printf("[totals of agency year] error getting data for first screen(estado:%s):%q", aID, err)
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro orgao=%s inválido", aID))
	}
	for _, agencyMonthlyInfo := range agenciesMonthlyInfo[aID] {
		if agencyMonthlyInfo.Summary != nil && agencyMonthlyInfo.Summary.BaseRemuneration.Total+agencyMonthlyInfo.Summary.OtherRemunerations.Total > 0 {
			monthTotals := models.MonthTotals{Month: agencyMonthlyInfo.Month,
				BaseRemuneration:   agencyMonthlyInfo.Summary.BaseRemuneration.Total,
				OtherRemunerations: agencyMonthlyInfo.Summary.OtherRemunerations.Total,
				CrawlingTimestamp:  agencyMonthlyInfo.CrawlingTimestamp,
			}
			monthTotalsOfYear = append(monthTotalsOfYear, monthTotals)

			// The status 4 is a report from crawlers that data is unavailable or malformed. By removing them from the API results, we make sure they are displayed as if there is no data.
			// Fazemos duas checagens no formato do ProcInfo para saber se ele é vazio pois alguns dados diferem, no banco de dados, quando o procinfo é nulo.
		} else if agencyMonthlyInfo.ProcInfo != nil && agencyMonthlyInfo.ProcInfo.String() != "" && agencyMonthlyInfo.ProcInfo.Status != 4 {
			monthTotals := models.MonthTotals{Month: agencyMonthlyInfo.Month,
				BaseRemuneration:   0,
				OtherRemunerations: 0,
				CrawlingTimestamp:  agencyMonthlyInfo.CrawlingTimestamp,
				Error:              &models.ProcError{Stdout: agencyMonthlyInfo.ProcInfo.Stdout, Stderr: agencyMonthlyInfo.ProcInfo.Stderr},
			}
			monthTotalsOfYear = append(monthTotalsOfYear, monthTotals)
		}
	}
	sort.Slice(monthTotalsOfYear, func(i, j int) bool {
		return monthTotalsOfYear[i].Month < monthTotalsOfYear[j].Month
	})
	pkg, _ := mgoClient.Db.GetPackage(strModels.PackageFilterOpts{AgencyID: &aID, Year: &year, Month: nil, Group: nil})
	agencyTotalsYear := models.AgencyTotalsYear{Year: year, Agency: agency, MonthTotals: monthTotalsOfYear, AgencyFullName: agency.Name, SummaryPackage: pkg}
	return c.JSON(http.StatusOK, agencyTotalsYear)
}

func getBasicInfoOfState(c echo.Context) error {
	yearOfConsult := time.Now().Year()
	stateName := c.Param("estado")
	agencies, err := pgClient.GetOPE(stateName, yearOfConsult)
	if err != nil {
		// That happens when there is no information on that year.
		log.Printf("[basic info state] first error getting data for first screen(ano:%d, estado:%s). Going to try again with last year:%q", yearOfConsult, stateName, err)
		yearOfConsult = yearOfConsult - 1

		agencies, err = pgClient.GetOPE(stateName, yearOfConsult)
		if err != nil {
			log.Printf("[basic info state] error getting data for first screen(ano:%d, estado:%s):%q", yearOfConsult, stateName, err)
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetros ano=%d ou estado=%s são inválidos", yearOfConsult, stateName))
		}
	}
	var agenciesBasic []models.AgencyBasic
	for k := range agencies {
		agenciesBasic = append(agenciesBasic, models.AgencyBasic{Name: agencies[k].ID, FullName: agencies[k].Name, AgencyCategory: agencies[k].Entity})
	}
	state := models.State{Name: stateName, ShortName: "", FlagURL: "", Agency: agenciesBasic}
	return c.JSON(http.StatusOK, state)
}

func getSalaryOfAgencyMonthYear(c echo.Context) error {
	month, err := strconv.Atoi(c.Param("mes"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro mês=%d", month))
	}
	year, err := strconv.Atoi(c.Param("ano"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d", year))
	}
	agencyName := strings.ToLower(c.Param("orgao"))
	agencyMonthlyInfo, _, err := pgClient.GetOMA(month, year, agencyName)
	if err != nil {
		log.Printf("[salary agency month year] error getting data for second screen(mes:%d ano:%d, orgao:%s):%q", month, year, agencyName, err)
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d, mês=%d ou nome do orgão=%s são inválidos", year, month, agencyName))
	}
	// Fazemos duas checagens no formato do ProcInfo para saber se ele é vazio pois alguns dados diferem, no banco de dados, quando o procinfo é nulo.
	if agencyMonthlyInfo.ProcInfo != nil && agencyMonthlyInfo.ProcInfo.String() != "" {
		var newEnv = agencyMonthlyInfo.ProcInfo.Env
		for _, omittedField := range conf.EnvOmittedFields {
			for i, field := range newEnv {
				if strings.Contains(field, omittedField) {
					newEnv[i] = omittedField + "= ##omitida##"
					break
				}
			}
		}
		agencyMonthlyInfo.ProcInfo.Env = newEnv
		return c.JSON(http.StatusPartialContent, models.ProcInfoResult{
			ProcInfo:          agencyMonthlyInfo.ProcInfo,
			CrawlingTimestamp: agencyMonthlyInfo.CrawlingTimestamp,
		})
	}
	return c.JSON(http.StatusOK, models.DataForChartAtAgencyScreen{
		Members:     agencyMonthlyInfo.Summary.IncomeHistogram,
		MaxSalary:   agencyMonthlyInfo.Summary.BaseRemuneration.Max,
		PackageURL:  agencyMonthlyInfo.Package.URL,
		PackageHash: agencyMonthlyInfo.Package.Hash,
		PackageSize: agencyMonthlyInfo.Package.Size,
	})
}

func getSummaryOfAgency(c echo.Context) error {
	year, err := strconv.Atoi(c.Param("ano"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d inválido", year))
	}
	month, err := strconv.Atoi(c.Param("mes"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro mês=%d", month))
	}
	agencyName := c.Param("orgao")
	agencyMonthlyInfo, agency, err := pgClient.GetOMA(month, year, agencyName)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d, mês=%d ou nome do orgão=%s são inválidos", year, month, agencyName))
	}
	agencySummary := models.AgencySummary{
		FullName:   agency.Name,
		TotalWage:  agencyMonthlyInfo.Summary.BaseRemuneration.Total,
		MaxWage:    agencyMonthlyInfo.Summary.BaseRemuneration.Max,
		TotalPerks: agencyMonthlyInfo.Summary.OtherRemunerations.Total,
		MaxPerk:    agencyMonthlyInfo.Summary.OtherRemunerations.Max,
		TotalRemuneration: agencyMonthlyInfo.Summary.BaseRemuneration.Total +
			agencyMonthlyInfo.Summary.OtherRemunerations.Total,
		TotalMembers: agencyMonthlyInfo.Summary.Count,
		CrawlingTime: agencyMonthlyInfo.CrawlingTimestamp,
		HasNext:      time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).In(loc).Before(time.Now().AddDate(0, 1, 0)),
		HasPrevious:  time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).In(loc).After(time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC).In(loc)),
	}
	return c.JSON(http.StatusOK, agencySummary)
}

func generalSummaryHandler(c echo.Context) error {
	agencies, err := pgClient.GetAgenciesCount()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Erro ao contar orgãos: %q", err))
	}
	collections, err := pgClient.GetNumberOfMonthsCollected()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Erro ao contar registros: %q", err))
	}
	fmonth, fyear, err := pgClient.Db.GetFirstDateWithMonthlyInfo()
	if err != nil {
		log.Printf("Error buscando dados - GetFirstDateWithRemunerationRecords: %q", err)
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Erro buscando primeiro registro de remuneração: %q", err))
	}
	fdate := time.Date(fyear, time.Month(fmonth), 2, 0, 0, 0, 0, time.UTC).In(loc)
	lmonth, lyear, err := pgClient.GetLastDateWithMonthlyInfo()
	if err != nil {
		log.Printf("Error buscando dados - GetLastDateWithRemunerationRecords: %q", err)
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Erro buscando último registro de remuneração: %q", err))
	}
	ldate := time.Date(lyear, time.Month(lmonth), 2, 0, 0, 0, 0, time.UTC).In(loc)
	remuValue, err := pgClient.Db.GetGeneralMonthlyInfo()
	if err != nil {
		log.Printf("Error buscando dados - GetGeneralRemunerationValue: %q", err)
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Erro buscando valor total de remuneração: %q", err))
	}
	return c.JSON(http.StatusOK, models.GeneralTotals{
		AgencyAmount:             int(agencies),
		MonthlyTotalsAmount:      int(collections),
		StartDate:                fdate,
		EndDate:                  ldate,
		RemunerationRecordsCount: int(collections),
		GeneralRemunerationValue: remuValue,
	})
}

func getGeneralRemunerationFromYear(c echo.Context) error {
	year, err := strconv.Atoi(c.Param("ano"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d inválido", year))
	}
	data, err := mgoClient.Db.GetGeneralMonthlyInfosFromYear(year)
	if err != nil {
		fmt.Println("Error searching for monthly info from year: %w", err)
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error buscando dados"))
	}
	return c.JSON(http.StatusOK, data)
}

func getAllAgencies(c echo.Context) error {
	agencies, err := pgClient.Db.GetAllAgencies()
	if err != nil {
		fmt.Println("Error while listing agencies: %w", err)
		return c.JSON(http.StatusInternalServerError, "Error while listing agencies")
	}
	return c.JSON(http.StatusOK, agencies)
}

func getAgencyById(c echo.Context) error {
	agencyName := c.Param("orgao")
	agency, err := pgClient.Db.GetAgency(agencyName)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Agency not found")
	}
	return c.JSON(http.StatusFound, agency)
}

func getMonthlyInfo(c echo.Context) error {
	year, err := strconv.Atoi(c.Param("ano"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d inválido", year))
	}
	agencyName := strings.ToLower(c.Param("orgao"))
	var monthlyInfo map[string][]strModels.AgencyMonthlyInfo
	month := c.Param("mes")
	if month != "" {
		m, err := strconv.Atoi(month)
		if err != nil {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro mes=%d inválido", m))
		}
		oma, _, err := pgClient.GetOMA(m, year, agencyName)
		if err != nil {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("Error getting OMA data"))
		}
		monthlyInfo = map[string][]strModels.AgencyMonthlyInfo{
			agencyName: {*oma},
		}
	} else {
		monthlyInfo, err = pgClient.Db.GetMonthlyInfo([]strModels.Agency{{ID: agencyName}}, year)
	}
	if err != nil {
		log.Printf("[totals of agency year] error getting data for first screen(ano:%d, estado:%s):%q", year, agencyName, err)
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d ou orgao=%s inválidos", year, agencyName))
	}

	if len(monthlyInfo[agencyName]) == 0 {
		return c.NoContent(http.StatusNotFound)
	}

	// we recreate all struct to controll the serialization better

	type Backup struct {
		URL  string `json:"url,omitempty"`
		Hash string `json:"hash,omitempty"`
		Size int64  `json:"size,omitempty"`
	}
	type DataSummary struct {
		Max     float64 `json:"max,omitempty"`
		Min     float64 `json:"min,omitempty"`
		Average float64 `json:"media,omitempty"`
		Total   float64 `json:"total,omitempty"`
	}
	type Summary struct {
		Count              int         `json:"quantidade,omitempty"`
		BaseRemuneration   DataSummary `json:"remuneracao_base,omitempty"`
		OtherRemunerations DataSummary `json:"outras_remuneracoes,omitempty"`
	}
	type Summaries struct {
		MemberActive Summary `json:"membros_ativos,omitempty"`
	}
	type Metadata struct {
		OpenFormat       bool   `json:"formato_aberto"`
		Access           string `json:"acesso,omitempty"`
		Extension        string `json:"extensao,omitempty"`
		StrictlyTabular  bool   `json:"dados_estritamente_tabulares,omitempty"`
		ConsistentFormat bool   `json:"manteve_consistencia_no_formato,omitempty"`
		HasEnrollment    bool   `json:"tem_matricula,omitempty"`
		HasCapacity      bool   `json:"tem_lotacao,omitempty"`
		HasPosition      bool   `json:"tem_cargo,omitempty"`
		BaseRevenue      string `json:"remuneracao_basica,omitempty"`
		OtherRecipes     string `json:"outras_receitas,omitempty"`
		Expenditure      string `json:"despesas,omitempty"`
	}
	type Score struct {
		Score             float64 `json:"indice_transparencia"`
		CompletenessScore float64 `json:"indice_completude"`
		EasinessScore     float64 `json:"indice_facilidade"`
	}
	type MIError struct {
		ErrorMessage string `json:"err_msg,omitempty"`
		Status       int32  `json:"status,omitempty"`
		Cmd          string `json:"cmd,omitempty"`
	}
	type SummaryzedMI struct {
		AgencyID string     `json:"id_orgao,omitempty"`
		Month    int        `json:"mes,omitempty"`
		Year     int        `json:"ano,omitempty"`
		Summary  *Summaries `json:"sumarios,omitempty"`
		Package  *Backup    `json:"pacote_de_dados,omitempty"`
		Meta     *Metadata  `json:"metadados,omitempty`
		Score    *Score     `json:"indice_transparencia,omitempty`
		Error    *MIError   `json:"error,omitempty"`
	}
	var summaryzedMI []SummaryzedMI
	for i := range monthlyInfo {
		for _, mi := range monthlyInfo[i] {
			// Fazemos duas checagens no formato do ProcInfo para saber se ele é vazio pois alguns dados diferem, no banco de dados, quando o procinfo é nulo.
			if mi.ProcInfo == nil || mi.ProcInfo.String() == "" {
				summaryzedMI = append(
					summaryzedMI,
					SummaryzedMI{
						AgencyID: mi.AgencyID,
						Error:    nil,
						Month:    mi.Month,
						Year:     mi.Year,
						Package: &Backup{
							URL:  formatDownloadUrl(mi.Package.URL),
							Hash: mi.Package.Hash,
							Size: mi.Package.Size,
						},
						Summary: &Summaries{
							MemberActive: Summary{
								Count: mi.Summary.Count,
								BaseRemuneration: DataSummary{
									Max:     mi.Summary.BaseRemuneration.Max,
									Min:     mi.Summary.BaseRemuneration.Min,
									Average: mi.Summary.BaseRemuneration.Average,
									Total:   mi.Summary.BaseRemuneration.Total,
								},
								OtherRemunerations: DataSummary{
									Max:     mi.Summary.OtherRemunerations.Max,
									Min:     mi.Summary.OtherRemunerations.Min,
									Average: mi.Summary.OtherRemunerations.Average,
									Total:   mi.Summary.OtherRemunerations.Total,
								},
							},
						},
						Meta: &Metadata{
							OpenFormat:       mi.Meta.OpenFormat,
							Access:           mi.Meta.Access,
							Extension:        mi.Meta.Extension,
							StrictlyTabular:  mi.Meta.StrictlyTabular,
							ConsistentFormat: mi.Meta.ConsistentFormat,
							HasEnrollment:    mi.Meta.HaveEnrollment,
							HasCapacity:      mi.Meta.ThereIsACapacity,
							HasPosition:      mi.Meta.HasPosition,
							BaseRevenue:      mi.Meta.BaseRevenue,
							OtherRecipes:     mi.Meta.OtherRecipes,
							Expenditure:      mi.Meta.Expenditure,
						},
						Score: &Score{
							Score:             mi.Score.Score,
							CompletenessScore: mi.Score.CompletenessScore,
							EasinessScore:     mi.Score.EasinessScore,
						}})
				// The status 4 is a report from crawlers that data is unavailable or malformed. By removing them from the API results, we make sure they are displayed as if there is no data.
			} else if mi.ProcInfo.Status != 4 {
				summaryzedMI = append(
					summaryzedMI,
					SummaryzedMI{
						AgencyID: mi.AgencyID,
						Error: &MIError{
							ErrorMessage: mi.ProcInfo.Stderr,
							Status:       mi.ProcInfo.Status,
							Cmd:          mi.ProcInfo.Cmd,
						},
						Month:   mi.Month,
						Year:    mi.Year,
						Package: nil,
						Summary: nil,
						Meta:    nil})
			}
		}
	}
	return c.JSON(http.StatusOK, summaryzedMI)
}

func searchByUrl(c echo.Context) error {
	//Pegando os query params
	years := c.QueryParam("anos")
	months := c.QueryParam("meses")
	agencies := c.QueryParam("orgaos")
	categories := c.QueryParam("categorias")
	types := c.QueryParam("tipos")
	//Criando os filtros a partir dos query params e validando eles
	filter, err := models.NewFilter(years, months, agencies, categories, types)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	var category string
	if filter != nil {
		category = filter.Category
	}
	// Pegando os resultados da pesquisa a partir dos filtros;
	results, err := postgresDb.Filter(postgresDb.RemunerationQuery(filter), postgresDb.Arguments(filter))
	if err != nil {
		log.Printf("Error querying BD (filter or counter):%q", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	remunerations, numRows, err := getSearchResults(conf.SearchLimit, category, results)
	if err != nil {
		log.Printf("Error getting search results: %q", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := models.SearchResponse{
		DownloadAvailable:  numRows > 0 && numRows <= conf.DownloadLimit,
		NumRowsIfAvailable: numRows,
		DownloadLimit:      conf.DownloadLimit,
		SearchLimit:        conf.SearchLimit,
		Results:            remunerations, // retornando os SearchLimit primeiros elementos.
	}
	return c.JSON(http.StatusOK, response)
}

func downloadByUrl(c echo.Context) error {
	//Pegando os query params
	years := c.QueryParam("anos")
	months := c.QueryParam("meses")
	agencies := c.QueryParam("orgaos")
	categories := c.QueryParam("categorias")
	types := c.QueryParam("tipos")

	//Criando os filtros a partir dos query params e validando eles
	filter, err := models.NewFilter(years, months, agencies, categories, types)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	results, err := postgresDb.Filter(postgresDb.RemunerationQuery(filter), postgresDb.Arguments(filter))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	searchResults, _, err := getSearchResults(conf.DownloadLimit, filter.Category, results)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.Response().Header().Set("Content-Disposition", "attachment; filename=dadosjusbr-remuneracoes.csv")
	c.Response().Header().Set("Content-Type", c.Response().Header().Get("Content-Type"))
	err = gocsv.Marshal(searchResults, c.Response().Writer)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("erro tentando fazer download do csv: %q", err))
	}
	return nil
}

func getSearchResults(limit int, category string, results []models.SearchDetails) ([]models.SearchResult, int, error) {
	searchResults := []models.SearchResult{}
	numRows := 0
	if len(results) == 0 {
		return searchResults, numRows, nil
	} else {
		// A razão para essa ordenação é que quando o usuário escolhe diversos órgãos
		// provavelmente ele prefere ver dados de todos eles. Dessa forma, aumentamos
		// as chances do preview limitado retornar dados de diversos órgãos.
		sort.SliceStable(results, func(i, j int) bool {
			return results[i].Ano < results[j].Ano || results[i].Mes < results[j].Mes
		})
		searchResults, numRows, err := sess.GetRemunerationsFromS3(limit, conf.DownloadLimit, category, conf.AwsS3Bucket, results)
		if err != nil {
			return nil, numRows, fmt.Errorf("failed to get remunerations from s3 %q", err)
		}
		return searchResults, numRows, nil
	}
}

func formatDownloadUrl(url string) string {
	return strings.Replace(url, conf.PackageRepoURL, conf.DadosJusURL, -1)
}

func main() {
	godotenv.Load() // There is no problem if the .env can not be loaded.
	l, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		log.Fatal(err.Error())
	}
	loc = l
	if err := envconfig.Process("api", &conf); err != nil {
		log.Fatal(err.Error())
	}

	// Criando o client do storage a partir do mongodb
	mongoDB, err := newMongoDB(conf)
	if err != nil {
		log.Fatal(err)
	}
	mgoClient, err = newClient(mongoDB)
	if err != nil {
		log.Fatal(err)
	}

	pgDB, err := newPostgresDB(conf)
	if err != nil {
		log.Fatal(err)
	}
	pgClient, err = newClient(pgDB)
	if err != nil {
		log.Fatal(err)
	}

	pgCredentials, err := NewPgCredentials(conf)
	if err != nil {
		log.Fatal("Error creating postgres credentials: %v", err)
	}

	postgresDb, err = NewPostgresDB(*pgCredentials)
	if err != nil {
		log.Fatalf("Error connecting to postgres: %v", err)
	}
	defer postgresDb.Disconnect()

	postgresDb.conn, err = pgDB.GetConnection()
	if err != nil {
		log.Fatalf("Error connecting to postgres: %v", err)
	}

	sess, err = NewAwsSession(conf.AwsRegion)
	if err != nil {
		log.Fatalf("Error creating aws session: %v", err)
	}

	fmt.Printf("Going to start listening at port:%d\n", conf.Port)

	e := echo.New()

	e.GET("/", func(ctx echo.Context) error { return nil }) // necessário para checagem do beanstalk.

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "ui/dist/",
		Browse: true,
		HTML5:  true,
		Index:  "index.html",
	}))
	e.Static("/static", "templates/assets")
	e.Use(middleware.Logger())

	// Internal API configuration
	uiAPIGroup := e.Group("/uiapi")
	if os.Getenv("DADOSJUSBR_ENV") == "Prod" {
		if conf.NewRelicApp == "" || conf.NewRelicLicense == "" {
			log.Fatalf("Missing environment variables NEWRELIC_APP_NAME or NEWRELIC_LICENSE")
		}
		nr, err := newrelic.NewApplication(
			newrelic.ConfigAppName(conf.NewRelicApp),
			newrelic.ConfigLicense(conf.NewRelicLicense),
			newrelic.ConfigAppLogForwardingEnabled(true),
		)
		if err != nil {
			log.Fatalf("Error bringin up new relic:%q", err)
		}
		postgresDb.newrelic = nr
		sess.newrelic = nr
		uiAPIGroup.Use(nrecho.Middleware(nr))
		uiAPIGroup.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"https://dadosjusbr.com", "http://dadosjusbr.com", "https://dadosjusbr.org", "http://dadosjusbr.org", "https://dadosjusbr-site-novo.herokuapp.com", "http://dadosjusbr-site-novo.herokuapp.com"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderContentLength},
		}))
		log.Println("Using production CORS")
	} else {
		uiAPIGroup.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderContentLength, echo.HeaderAccessControlAllowOrigin},
		}))
	}
	// Return a summary of an agency. This information will be used in the head of the agency page.
	uiAPIGroup.GET("/v1/orgao/resumo/:orgao/:ano/:mes", getSummaryOfAgency)
	// Return all the salary of a month and year. This will be used in the point chart at the entity page.
	uiAPIGroup.GET("/v1/orgao/salario/:orgao/:ano/:mes", getSalaryOfAgencyMonthYear)
	// Return the total of salary of every month of a year of a agency. The salary is divided in Wage, Perks and Others. This will be used to plot the bars chart at the state page.
	uiAPIGroup.GET("/v1/orgao/totais/:orgao/:ano", getTotalsOfAgencyYear)
	// Return basic information of a state
	uiAPIGroup.GET("/v1/orgao/:estado", getBasicInfoOfState)
	uiAPIGroup.GET("/v1/geral/remuneracao/:ano", getGeneralRemunerationFromYear)
	uiAPIGroup.GET("/v1/geral/resumo", generalSummaryHandler)
	// Retorna um conjunto de dados a partir de filtros informados por query params
	uiAPIGroup.GET("/v2/pesquisar", searchByUrl)
	// Baixa um conjunto de dados a partir de filtros informados por query params
	uiAPIGroup.GET("/v2/download", downloadByUrl)

	// Public API configuration
	apiGroup := e.Group("/v1", middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderContentLength},
	}))
	// Return agency
	apiGroup.GET("/orgao/:orgao", getAgencyById)
	// Return all agencies
	apiGroup.GET("/orgaos", getAllAgencies)
	// Return MIs by year
	apiGroup.GET("/dados/:orgao/:ano", getMonthlyInfo)
	// Return MIs by month
	apiGroup.GET("/dados/:orgao/:ano/:mes", getMonthlyInfo)

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", conf.Port),
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	}
	e.Logger.Fatal(e.StartServer(s))
}
