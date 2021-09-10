package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dadosjusbr/remuneracao-magistrados/models"
	"github.com/dadosjusbr/storage"
	"github.com/joho/godotenv"

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

	// Omited fields
	EnvOmittedFields []string `envconfig:"ENV_OMITTED_FIELDS"`
}

var client *storage.Client
var loc *time.Location

// newClient takes a config struct and creates a client to connect with DB and Cloud5
func newClient(c config) (*storage.Client, error) {
	if c.MongoMICol == "" || c.MongoAgCol == "" {
		return nil, fmt.Errorf("error creating storage client: db collections must not be empty. MI:\"%s\", AG:\"%s\", PKG:\"%s\"", c.MongoMICol, c.MongoAgCol, c.MongoPkgCol)
	}
	db, err := storage.NewDBClient(c.MongoURI, c.MongoDBName, c.MongoMICol, c.MongoAgCol, c.MongoPkgCol)
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
		if agencyMonthlyInfo.Summary.MemberActive.Wage.Total+agencyMonthlyInfo.Summary.MemberActive.Perks.Total+agencyMonthlyInfo.Summary.MemberActive.Others.Total > 0 {
			monthTotals := models.MonthTotals{Month: agencyMonthlyInfo.Month,
				Wage:   agencyMonthlyInfo.Summary.MemberActive.Wage.Total,
				Perks:  agencyMonthlyInfo.Summary.MemberActive.Perks.Total,
				Others: agencyMonthlyInfo.Summary.MemberActive.Others.Total,
			}
			monthTotalsOfYear = append(monthTotalsOfYear, monthTotals)
		}
	}
	sort.Slice(monthTotalsOfYear, func(i, j int) bool {
		return monthTotalsOfYear[i].Month < monthTotalsOfYear[j].Month
	})
	pkg, _ := client.Db.GetPackage(storage.PackageFilterOpts{AgencyID: &aID, Year: &year, Month: nil, Group: nil})
	agencyTotalsYear := models.AgencyTotalsYear{Year: year, MonthTotals: monthTotalsOfYear, AgencyFullName: agency.Name, SummaryPackage: pkg}
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
		Members:     agencyMonthlyInfo.Summary.MemberActive.IncomeHistogram,
		MaxSalary:   agencyMonthlyInfo.Summary.MemberActive.Wage.Max,
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
		FullName:  agency.Name,
		TotalWage: agencyMonthlyInfo.Summary.MemberActive.Wage.Total,
		MaxWage:   agencyMonthlyInfo.Summary.MemberActive.Wage.Max,
		TotalPerks: agencyMonthlyInfo.Summary.MemberActive.Perks.Total +
			agencyMonthlyInfo.Summary.MemberActive.Others.Total,
		MaxPerk: math.Max(agencyMonthlyInfo.Summary.MemberActive.Perks.Max, agencyMonthlyInfo.Summary.MemberActive.Others.Max),
		TotalRemuneration: agencyMonthlyInfo.Summary.MemberActive.Wage.Total +
			agencyMonthlyInfo.Summary.MemberActive.Perks.Total +
			agencyMonthlyInfo.Summary.MemberActive.Others.Total +
			agencyMonthlyInfo.Summary.MemberActive.Wage.Total,
		TotalMembers: agencyMonthlyInfo.Summary.MemberActive.Count,
		CrawlingTime: agencyMonthlyInfo.CrawlingTimestamp,
		HasNext:      verifyNextOMA(month, year, agencyName),
		HasPrevious:  verifyPreviousOMA(month, year, agencyName),
	}
	return c.JSON(http.StatusOK, agencySummary)
}

func verifyNextOMA(month int, year int, agencyName string) bool {
	if month == 12 {
		month = 1
		year += 1
	} else {
		month += 1
	}
	_, _, err := client.GetOMA(month, year, agencyName)
	return err == nil
}
func verifyPreviousOMA(month int, year int, agencyName string) bool {
	if month == 1 {
		month = 12
		year -= 1
	} else {
		month -= 1
	}
	_, _, err := client.GetOMA(month, year, agencyName)
	return err == nil
}

func apiOMA(c echo.Context) error {
	year, err := strconv.Atoi(c.Param("ano"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d inválido", year))
	}
	month, err := strconv.Atoi(c.Param("mes"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro mês=%d", month))
	}
	agName := c.Param("orgao")
	agMI, _, err := client.GetOMA(month, year, agName)
	if err != nil {
		c.Logger().Printf("Error fetching data for API (%s?%s):%q", c.Path(), c.QueryString(), err)
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error buscando dados"))
	}
	switch format := c.QueryParam("format"); format {
	case "zip":
		return c.Redirect(http.StatusPermanentRedirect, agMI.Package.URL)
	case "json", "":
		return c.JSONPretty(http.StatusOK, agMI.Employee, " ")
	default:
		return c.String(http.StatusBadRequest, fmt.Sprintf("Por favor, escolher o formato entre json e zip!"))
	}
}

func generalSummaryHandler(c echo.Context) error {
	agencyAmount, err := client.GetAgenciesCount()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error buscando dados"))
	}
	miCount, err := client.GetNumberOfMonthsCollected()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error buscando dados"))
	}
	fmonth, fyear, err := client.GetFirstDateWithMonthlyInfo()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error buscando dados"))
	}
	lmonth, lyear, err := client.GetLastDateWithMonthlyInfo()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error buscando dados"))
	}
	remunerationSummary, err := client.Db.GetRemunerationSummary()
	if err != nil {
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
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error buscando dados"))
	}
	return c.JSON(http.StatusOK, data)
}

func getAllAgencies(c echo.Context) error {
	agencies, err := client.Db.GetAllAgencies()
	if err != nil {
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
	monthlyInfo, err := client.Db.GetMonthlyInfo([]storage.Agency{{ID: agencyName}}, year)
	if err != nil {
		log.Printf("[totals of agency year] error getting data for first screen(ano:%d, estado:%s):%q", year, agencyName, err)
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d ou orgao=%s inválidos", year, agencyName))
	}
	return c.JSON(http.StatusOK, monthlyInfo)
}

var conf config

func main() {
	godotenv.Load() // There is no problem if the .env can not be loaded.
	l, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		log.Fatal(err.Error())
	}
	loc = l
	if err := envconfig.Process("remuneracao-magistrados", &conf); err != nil {
		log.Fatal(err.Error())
	}

	// Criando o client do storage
	client, err = newClient(conf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Going to start listening at port:%d\n", conf.Port)

	e := echo.New()

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "ui/dist/",
		Browse: true,
		HTML5:  true,
		Index:  "index.html",
	}))
	e.Static("/static", "templates/assets")

	// Internal API configuration
	uiAPIGroup := e.Group("/uiapi")
	if os.Getenv("DADOSJUSBR_ENV") == "Prod" {
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

	// Public API configuration
	apiGroup := e.Group("/", middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderContentLength},
	}))
	// Return OMA (órgão/mês/ano) information
	apiGroup.GET("v1/orgao/:orgao/:ano/:mes", apiOMA)

	apiGroup.GET("v1/orgao/:orgao", getAgencyById)
	apiGroup.GET("v1/orgaos", getAllAgencies)
	apiGroup.GET("v1/resumo/:orgao/:ano", getMonthlyInfo)

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", conf.Port),
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	}
	e.Logger.Fatal(e.StartServer(s))
}
