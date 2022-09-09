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

var client *storage.Client
var loc *time.Location
var conf config
var postgresDb *PostgresDB

// newClient takes a config struct and creates a client to connect with DB and Cloud5
func newClient(c config) (*storage.Client, error) {
	if c.MongoMICol == "" || c.MongoAgCol == "" {
		return nil, fmt.Errorf("error creating storage client: db collections must not be empty. MI:\"%s\", AG:\"%s\", PKG:\"%s\"", c.MongoMICol, c.MongoAgCol, c.MongoPkgCol)
	}
	db, err := storage.NewDBClient(c.MongoURI, c.MongoDBName, c.MongoMICol, c.MongoAgCol, c.MongoPkgCol, c.MongoRevCol)
	if err != nil {
		return nil, fmt.Errorf("error creating DB client: %q", err)
	}
	db.Collection(c.MongoMICol)
	client, err := storage.NewClient(db, &storage.CloudClient{})
	if err != nil {
		return nil, fmt.Errorf("error creating storage.client: %q", err)
	}
	return client, nil
}

func getTotalsOfAgencyYear(c echo.Context) error {
	year, err := strconv.Atoi(c.Param("ano"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d inválido", year))
	}
	aID := c.Param("orgao")
	agenciesMonthlyInfo, err := client.Db.GetMonthlyInfo([]storage.Agency{{ID: aID}}, year)
	if err != nil {
		log.Printf("[totals of agency year] error getting data for first screen(ano:%d, estado:%s):%q", year, aID, err)
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d ou orgao=%s inválidos", year, aID))
	}
	var monthTotalsOfYear []models.MonthTotals
	agency, err := client.Db.GetAgency(aID)
	if err != nil {
		log.Printf("[totals of agency year] error getting data for first screen(estado:%s):%q", aID, err)
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro orgao=%s inválido", aID))
	}
	for _, agencyMonthlyInfo := range agenciesMonthlyInfo[aID] {
		if agencyMonthlyInfo.Summary.BaseRemuneration.Total+agencyMonthlyInfo.Summary.OtherRemunerations.Total > 0 {
			monthTotals := models.MonthTotals{Month: agencyMonthlyInfo.Month,
				BaseRemuneration:   agencyMonthlyInfo.Summary.BaseRemuneration.Total,
				OtherRemunerations: agencyMonthlyInfo.Summary.OtherRemunerations.Total,
				CrawlingTimestamp:  agencyMonthlyInfo.CrawlingTimestamp,
			}
			monthTotalsOfYear = append(monthTotalsOfYear, monthTotals)

			// The status 4 is a report from crawlers that data is unavailable or malformed. By removing them from the API results, we make sure they are displayed as if there is no data.
		} else if agencyMonthlyInfo.ProcInfo != nil && agencyMonthlyInfo.ProcInfo.Status != 4 {
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
	pkg, _ := client.Db.GetPackage(storage.PackageFilterOpts{AgencyID: &aID, Year: &year, Month: nil, Group: nil})
	agencyTotalsYear := models.AgencyTotalsYear{Year: year, Agency: agency, MonthTotals: monthTotalsOfYear, AgencyFullName: agency.Name, SummaryPackage: pkg}
	return c.JSON(http.StatusOK, agencyTotalsYear)
}

func getBasicInfoOfState(c echo.Context) error {
	yearOfConsult := time.Now().Year()
	stateName := c.Param("estado")
	agencies, _, err := client.GetOPE(stateName, yearOfConsult)
	if err != nil {
		// That happens when there is no information on that year.
		log.Printf("[basic info state] first error getting data for first screen(ano:%d, estado:%s). Going to try again with last year:%q", yearOfConsult, stateName, err)
		yearOfConsult = yearOfConsult - 1

		agencies, _, err = client.GetOPE(stateName, yearOfConsult)
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
	agencyName := c.Param("orgao")
	agencyMonthlyInfo, _, err := client.GetOMA(month, year, agencyName)
	if err != nil {
		log.Printf("[salary agency month year] error getting data for second screen(mes:%d ano:%d, orgao:%s):%q", month, year, agencyName, err)
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d, mês=%d ou nome do orgão=%s são inválidos", year, month, agencyName))
	}

	if agencyMonthlyInfo.ProcInfo != nil {
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
		PackageHash: agencyMonthlyInfo.Package.URL,
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
	agencyMonthlyInfo, agency, err := client.GetOMA(month, year, agencyName)
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
	agencyAmount, err := client.GetAgenciesCount()
	if err != nil {
		log.Printf("Error buscando dados - GetAgenciesCount: %q", err)
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error buscando dados:"))
	}
	miCount, err := client.GetNumberOfMonthsCollected()
	if err != nil {
		log.Printf("Error buscando dados - GetNumberOfMonthsCollected: %q", err)
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error buscando dados"))
	}
	fmonth, fyear, err := client.GetFirstDateWithMonthlyInfo()
	if err != nil {
		log.Printf("Error buscando dados - GetFirstDateWithMonthlyInfo: %q", err)
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error buscando dados"))
	}
	lmonth, lyear, err := client.GetLastDateWithMonthlyInfo()
	if err != nil {
		log.Printf("Error buscando dados - GetLastDateWithMonthlyInfo: %q", err)
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error buscando dados"))
	}
	remunerationSummary, err := client.Db.GetRemunerationSummary()
	if err != nil {
		log.Printf("Error buscando dados - GetRemunerationSummary: %q", err)
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error buscando dados"))
	}
	fdate := time.Date(fyear, time.Month(fmonth), 2, 0, 0, 0, 0, time.UTC).In(loc)
	ldate := time.Date(lyear, time.Month(lmonth), 2, 0, 0, 0, 0, time.UTC).In(loc)
	return c.JSON(http.StatusOK, models.GeneralTotals{
		AgencyAmount:             agencyAmount,
		MonthlyTotalsAmount:      miCount,
		StartDate:                fdate,
		EndDate:                  ldate,
		RemunerationRecordsCount: remunerationSummary.Count,
		GeneralRemunerationValue: remunerationSummary.Value})
}

func getGeneralRemunerationFromYear(c echo.Context) error {
	year, err := strconv.Atoi(c.Param("ano"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d inválido", year))
	}
	data, err := client.Db.GetGeneralMonthlyInfosFromYear(year)
	if err != nil {
		fmt.Println("Error searching for monthly info from year: %w", err)
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error buscando dados"))
	}
	return c.JSON(http.StatusOK, data)
}

func getAllAgencies(c echo.Context) error {
	agencies, err := client.Db.GetAllAgencies()
	if err != nil {
		fmt.Println("Error while listing agencies: %w", err)
		return c.JSON(http.StatusInternalServerError, "Error while listing agencies")
	}
	return c.JSON(http.StatusOK, agencies)
}

func getAgencyById(c echo.Context) error {
	agencyName := c.Param("orgao")
	agency, err := client.Db.GetAgency(agencyName)
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
	agencyName := c.Param("orgao")
	var monthlyInfo map[string][]storage.AgencyMonthlyInfo
	month := c.Param("mes")
	if month != "" {
		m, err := strconv.Atoi(month)
		if err != nil {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro mes=%d inválido", m))
		}
		oma, _, err := client.GetOMA(m, year, agencyName)
		if err != nil {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("Error getting OMA data"))
		}
		monthlyInfo = map[string][]storage.AgencyMonthlyInfo{
			agencyName: {*oma},
		}
	} else {
		monthlyInfo, err = client.Db.GetMonthlyInfo([]storage.Agency{{ID: agencyName}}, year)
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
		NoLoginRequired   bool   `json:"login_nao_necessario"`
		NoCaptchaRequired bool   `json:"captcha_nao_necessario"`
		Access            string `json:"acesso,omitempty"`
		Extension         string `json:"extensao,omitempty"`
		StrictlyTabular   bool   `json:"dados_estritamente_tabulares,omitempty"`
		ConsistentFormat  bool   `json:"manteve_consistencia_no_formato,omitempty"`
		HasEnrollment     bool   `json:"tem_matricula,omitempty"`
		HasCapacity       bool   `json:"tem_lotacao,omitempty"`
		HasPosition       bool   `json:"tem_cargo,omitempty"`
		BaseRevenue       string `json:"remuneracao_basica,omitempty"`
		OtherRecipes      string `json:"outras_receitas,omitempty"`
		Expenditure       string `json:"despesas,omitempty"`
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
			if mi.ProcInfo == nil {
				summaryzedMI = append(summaryzedMI, SummaryzedMI{AgencyID: mi.AgencyID, Error: nil, Month: mi.Month, Year: mi.Year, Package: &Backup{
					URL:  formatDownloadUrl(mi.Package.URL),
					Hash: mi.Package.Hash,
				}, Summary: &Summaries{
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
				}, Meta: &Metadata{
					NoLoginRequired:   mi.Meta.NoLoginRequired,
					NoCaptchaRequired: mi.Meta.NoCaptchaRequired,
					Access:            mi.Meta.Access,
					Extension:         mi.Meta.Extension,
					StrictlyTabular:   mi.Meta.StrictlyTabular,
					ConsistentFormat:  mi.Meta.ConsistentFormat,
					HasEnrollment:     mi.Meta.HaveEnrollment,
					HasCapacity:       mi.Meta.ThereIsACapacity,
					HasPosition:       mi.Meta.HasPosition,
					BaseRevenue:       mi.Meta.BaseRevenue,
					OtherRecipes:      mi.Meta.OtherRecipes,
					Expenditure:       mi.Meta.Expenditure,
				}, Score: &Score{
					Score:             mi.Score.Score,
					CompletenessScore: mi.Score.CompletenessScore,
					EasinessScore:     mi.Score.EasinessScore,
				}})
				// The status 4 is a report from crawlers that data is unavailable or malformed. By removing them from the API results, we make sure they are displayed as if there is no data.
			} else if mi.ProcInfo.Status != 4 {
				summaryzedMI = append(summaryzedMI, SummaryzedMI{AgencyID: mi.AgencyID, Error: &MIError{
					ErrorMessage: mi.ProcInfo.Stderr,
					Status:       mi.ProcInfo.Status,
					Cmd:          mi.ProcInfo.Cmd,
				}, Month: mi.Month, Year: mi.Year, Package: nil, Summary: nil, Meta: nil})
			}
		}
	}
	return c.JSON(http.StatusOK, summaryzedMI)
}

func searchByUrl(c echo.Context) error {
	var years string
	var months string
	var agencies string
	var categories string
	var types string
	//Pegando os query params
	years = c.QueryParam("anos")
	months = c.QueryParam("meses")
	agencies = c.QueryParam("orgaos")
	categories = c.QueryParam("categorias")
	types = c.QueryParam("tipos")

	//Criando os filtros a partir dos query params e validando eles
	filter, err := models.NewFilter(years, months, agencies, categories, types)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Pegando os resultados da pesquisa a partir dos filtros;
	results, err := postgresDb.Filter(remunerationQuery(filter, conf.DownloadLimit+1), arguments(filter))
	if err != nil {
		log.Printf("Error querying BD (filter or counter):%q", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	returnedResults := []models.SearchResult{}
	// Devemos retornar um array vazio quando a pesquisa não retornar dados.
	if len(results) > 0 {
		// Nesse caso, precisamos checamos se a quantidade de resultados
		// é menor que o search limit (para evitar array out of bounds)
		upper := conf.SearchLimit
		if len(results) < conf.SearchLimit {
			upper = len(results)
		}
		returnedResults = results[0:upper]
	}
	response := models.SearchResponse{
		DownloadAvailable:  len(results) > 0 && len(results) <= conf.DownloadLimit,
		NumRowsIfAvailable: len(results),
		DownloadLimit:      conf.DownloadLimit,
		SearchLimit:        conf.SearchLimit,
		Results:            returnedResults, // retornando os SearchLimit primeiros elementos.
	}
	return c.JSON(http.StatusOK, response)
}

func downloadByUrl(c echo.Context) error {
	var years string
	var months string
	var agencies string
	var categories string
	var types string
	//Pegando os query params
	years = c.QueryParam("anos")
	months = c.QueryParam("meses")
	agencies = c.QueryParam("orgaos")
	categories = c.QueryParam("categorias")
	types = c.QueryParam("tipos")

	//Criando os filtros a partir dos query params e validando eles
	filter, err := models.NewFilter(years, months, agencies, categories, types)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	results, err := postgresDb.Filter(remunerationQuery(filter, conf.DownloadLimit), arguments(filter))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.Response().Header().Set("Content-Disposition", "attachment; filename=dadosjusbr-remuneracoes.csv")
	c.Response().Header().Set("Content-Type", c.Response().Header().Get("Content-Type"))
	err = gocsv.Marshal(results, c.Response().Writer)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("erro tentando fazer download do csv: %q", err))
	}
	return nil
}

//Função que recebe os filtros e a partir deles estrutura a query SQL da pesquisa
func remunerationQuery(filter *models.Filter, limit int) string {
	//A query padrão sem os filtros
	query := ` 
	SELECT 
		c.id_orgao as orgao,
		c.mes as mes,
		c.ano as ano,
		c.matricula AS matricula,
		c.nome AS nome, 
		c.cargo as cargo,
		c.lotacao as lotacao,
		r.categoria_contracheque as categoria_contracheque,
		r.detalhamento_contracheque as detalhamento_contracheque,
		r.valor as valor 
	FROM contracheques c
		INNER JOIN remuneracoes r ON r.id_coleta = c.id_coleta AND r.id_contracheque = c.id
	`
	if filter != nil {
		addFiltersInQuery(&query, filter)
	}

	return fmt.Sprintf("%s FETCH FIRST %d ROWS ONLY;", query, limit)
}

//Função que define os argumentos passados para a query
func arguments(filter *models.Filter) []interface{} {
	var arguments []interface{}
	if filter != nil {
		if len(filter.Years) > 0 {
			for _, y := range filter.Years {
				arguments = append(arguments, y)
			}
		}
		if len(filter.Months) > 0 {
			for _, m := range filter.Months {
				arguments = append(arguments, m)
			}
		}
		if len(filter.Agencies) > 0 {
			for _, a := range filter.Agencies {
				arguments = append(arguments, a)
			}
		}
		if len(filter.Categories) > 0 {
			for _, c := range filter.Categories {
				arguments = append(arguments, c)
			}
		}
		if filter.Types != "" {
			// Adicionando '% %' na clausura LIKE
			arguments = append(arguments, fmt.Sprintf("%%%s%%", filter.Types))
		}
	}

	return arguments
}

//Função que insere os filtros na query
func addFiltersInQuery(query *string, filter *models.Filter) {
	*query = *query + " WHERE"

	//Insere os filtros de ano caso existam
	if len(filter.Years) > 0 {
		for i := 0; i < len(filter.Years); i++ {
			if i == 0 {
				*query = fmt.Sprintf("%s (", *query)
			}
			*query = fmt.Sprintf("%s c.ano = $%d", *query, i+1)
			if i < len(filter.Years)-1 {
				*query = fmt.Sprintf("%s OR", *query)
			}
		}
		*query = fmt.Sprintf("%s)", *query)
	}

	//Insere os filtros de mês
	if len(filter.Months) > 0 {
		lastIndex := len(filter.Years)
		if lastIndex > 0 {
			*query = fmt.Sprintf("%s AND", *query)
		}
		for i := lastIndex; i < len(filter.Months)+lastIndex; i++ {
			if i == lastIndex {
				*query = fmt.Sprintf("%s (", *query)
			}
			*query = fmt.Sprintf("%s c.mes = $%d", *query, i+1)
			if i < len(filter.Months)+lastIndex-1 {
				*query = fmt.Sprintf("%s OR", *query)
			}
		}
		*query = fmt.Sprintf("%s)", *query)
	}

	//Insere o filtro de órgãos
	if len(filter.Agencies) > 0 {
		lastIndex := len(filter.Years) + len(filter.Months)
		if lastIndex > 0 {
			*query = fmt.Sprintf("%s AND", *query)
		}
		for i := lastIndex; i < lastIndex+len(filter.Agencies); i++ {
			if i == lastIndex {
				*query = fmt.Sprintf("%s (", *query)
			}
			*query = fmt.Sprintf("%s c.id_orgao = $%d", *query, i+1)
			if i < lastIndex+len(filter.Agencies)-1 {
				*query = fmt.Sprintf("%s OR", *query)
			}
		}
		*query = fmt.Sprintf("%s)", *query)
	}

	//Insere o filtro de categoria das remunerações
	if len(filter.Categories) > 0 {
		lastIndex := len(filter.Years) + len(filter.Months) + len(filter.Agencies)
		if lastIndex > 0 {
			*query = fmt.Sprintf("%s AND", *query)
		}
		for i := lastIndex; i < lastIndex+len(filter.Categories); i++ {
			if i == lastIndex {
				*query = fmt.Sprintf("%s (", *query)
			}
			*query = fmt.Sprintf("%s r.categoria_contracheque = $%d", *query, i+1)
			if i < lastIndex+len(filter.Categories)-1 {
				*query = fmt.Sprintf("%s OR", *query)
			}
		}
		*query = fmt.Sprintf("%s)", *query)
	}

	//Insere o filtro do tipo de órgãos
	if filter.Types != "" {
		lastIndex := len(filter.Years) + len(filter.Months) + len(filter.Agencies) + len(filter.Categories)
		*query = fmt.Sprintf("%s AND ( c.id_orgao like $%d )", *query, lastIndex+1)
	}
}

func formatDownloadUrl(url string) string {
	return strings.Replace(url, conf.PackageRepoURL, conf.DadosJusURL, -1)
}

func countRemunerationQuery(filter *models.Filter) string {
	query := ` 
	SELECT 
		COUNT(*)
	FROM contracheques c
		INNER JOIN remuneracoes r ON r.id_coleta = c.id_coleta AND r.id_contracheque = c.id
	`
	if filter != nil {
		addFiltersInQuery(&query, filter)
	}
	return fmt.Sprintf("%s;", query)
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

	// Criando o client do storage
	client, err = newClient(conf)
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

	// Internal API configuration
	uiAPIGroup := e.Group("/uiapi")
	uiAPIGroup.Use(middleware.Logger())
	if os.Getenv("DADOSJUSBR_ENV") == "Prod" {
		if conf.NewRelicApp == "" || conf.NewRelicLicense == "" {
			log.Fatalf("Missing environment variables NEWRELIC_APP_NAME or NEWRELIC_LICENSE")
		}
		nr, err := newrelic.NewApplication(
			newrelic.ConfigAppName(conf.NewRelicApp),
			newrelic.ConfigLicense(conf.NewRelicLicense),
			newrelic.ConfigAppLogForwardingEnabled(true),
		)
		postgresDb.newrelic = nr
		if err != nil {
			log.Fatalf("Error bringin up new relic:%q", err)
		}
		uiAPIGroup.Use(nrecho.Middleware(nr))
		uiAPIGroup.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"https://dadosjusbr.com", "http://dadosjusbr.com", "https://dadosjusbr.org", "http://dadosjusbr.org", "https://dadosjusbr-site-novo.herokuapp.com", "http://dadosjusbr-site-novo.herokuapp.com"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderContentLength},
		}))
		log.Println("Using production CORS")
	} else {
		uiAPIGroup.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderContentLength},
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
