package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dadosjusbr/remuneracao-magistrados/db"
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

	// Omited fields
	EnvOmittedFields []string `envconfig:"ENV_OMITTED_FIELDS"`
}

var monthsLabelMap = map[int]string{
	1:  "Janeiro",
	2:  "Fevereiro",
	3:  "Março",
	4:  "Abril",
	5:  "Maio",
	6:  "Junho",
	7:  "Julho",
	8:  "Agosto",
	9:  "Setembro",
	10: "Outubro",
	11: "Novembro",
	12: "Dezembro",
}

var client *storage.Client

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

//SidebarElement contains the necessary info to render the sidebar
type SidebarElement struct {
	Label string
	URL   string
}

func getSidebarElements(dbClient *db.Client) ([]SidebarElement, error) {
	processedMonths, err := dbClient.GetProcessedMonths()
	if err != nil {
		return nil, fmt.Errorf("error retrieving all processed months from db --> %v", err)
	}

	var sidebarElements []SidebarElement

	for _, pm := range processedMonths {
		label := fmt.Sprintf("%s %d", monthsLabelMap[pm.Month], pm.Year)
		URL := fmt.Sprintf("/%d/%d", pm.Year, pm.Month)
		sidebarElements = append(sidebarElements, SidebarElement{Label: label, URL: URL})
	}

	return sidebarElements, nil
}

func handleMonthRequest(dbClient *db.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		month, err := strconv.Atoi(c.Param("month"))
		if err != nil {
			fmt.Println(fmt.Errorf("invalid month on the url: (%s) --> %v", c.Param("month"), err))
			return c.String(http.StatusBadRequest, "invalid month")
		}
		year, err := strconv.Atoi(c.Param("year"))
		if err != nil {
			fmt.Println(fmt.Errorf("invalid year on the url: (%s) --> %v", c.Param("year"), err))
			return c.String(http.StatusBadRequest, "invalid year")
		}

		monthResults, err := dbClient.GetMonthResults(month, year)
		if err != nil {
			if err == db.ErrDocNotFound {
				//TODO: render a 404 page
				fmt.Println("Document not found")
				return c.String(http.StatusNotFound, "not found")
			}
			fmt.Println(fmt.Errorf("unexpected error fetching month data from DB --> %v", err))
			return c.String(http.StatusInternalServerError, "unexpected error")
		}

		monthLabel := fmt.Sprintf("%s %d", monthsLabelMap[month], year)

		sidebarElements, err := getSidebarElements(dbClient)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusInternalServerError, "unexpected error")
		}

		viewModel := struct {
			Month           int
			Year            int
			MonthLabel      string
			SpreadsheetsURL string
			DatapackageURL  string
			SidebarElements []SidebarElement
			Statistics      []db.Statistic
		}{
			monthResults.Month,
			monthResults.Year,
			monthLabel,
			monthResults.SpreadsheetsURL,
			monthResults.DatapackageURL,
			sidebarElements,
			monthResults.Statistics,
		}
		return c.Render(http.StatusOK, "monthTemplate.html", viewModel)
	}
}

func handleMainPageRequest(dbClient *db.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		sidebarElements, err := getSidebarElements(dbClient)
		if err != nil {
			log.Printf("Error getting sidebar elements: %q", err)
			return c.String(http.StatusInternalServerError, "unexpected error")
		}
		viewModel := struct {
			SidebarElements []SidebarElement
		}{
			sidebarElements,
		}
		return c.Render(http.StatusOK, "homePageTemplate.html", viewModel)
	}
}

// newClient takes a config struct and creates a client to connect with DB and Cloud5
func newClient(c config) (*storage.Client, error) {
	if c.MongoMICol == "" || c.MongoAgCol == "" {
		return nil, fmt.Errorf("error creating storage client: db collections must not be empty. MI:\"%s\", AG:\"%s\"", c.MongoMICol, c.MongoAgCol)
	}
	db, err := storage.NewDBClient(c.MongoURI, c.MongoDBName, c.MongoMICol, c.MongoAgCol)
	if err != nil {
		return nil, fmt.Errorf("error creating DB client: %q", err)
	}
	db.Collection(c.MongoMICol)
	client, err := storage.NewClient(db, &storage.BackupClient{})
	if err != nil {
		return nil, fmt.Errorf("error creating storage.client: %q", err)
	}
	return client, nil
}

func getTotalsOfAgencyYear(c echo.Context) error {
	stateName := c.Param("estado")
	year, err := strconv.Atoi(c.Param("ano"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d inválido", year))
	}

	_, agenciesMonthlyInfo, err := client.GetDataForFirstScreen(stateName, year)
	if err != nil {
		log.Printf("[totals of agency year] error getting data for first screen(ano:%d, estado:%s):%q", year, stateName, err)
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d ou estado=%s inválidos", year, stateName))
	}
	var monthTotalsOfYear []models.MonthTotals
	agencyName := c.Param("orgao")
	for _, agencyMonthlyInfo := range agenciesMonthlyInfo[agencyName] {
		monthTotals := models.MonthTotals{Month: agencyMonthlyInfo.Month,
			Wage:   agencyMonthlyInfo.Summary.General.Wage.Total,
			Perks:  agencyMonthlyInfo.Summary.General.Perks.Total,
			Others: agencyMonthlyInfo.Summary.General.Others.Total,
		}
		monthTotalsOfYear = append(monthTotalsOfYear, monthTotals)
	}
	sort.Slice(monthTotalsOfYear, func(i, j int) bool {
		return monthTotalsOfYear[i].Month < monthTotalsOfYear[j].Month
	})
	agencyTotalsYear := models.AgencyTotalsYear{Year: year, MonthTotals: monthTotalsOfYear}
	return c.JSON(http.StatusOK, agencyTotalsYear)
}

func getBasicInfoOfState(c echo.Context) error {
	yearOfConsult := time.Now().Year()
	stateName := c.Param("estado")
	agencies, _, err := client.GetDataForFirstScreen(stateName, yearOfConsult)
	if err != nil {
		log.Printf("[basic info state] first error getting data for first screen(ano:%d, estado:%s). Going to try again with last year:%q", yearOfConsult, stateName, err)
		yearOfConsult = yearOfConsult - 1
	}
	agencies, _, err = client.GetDataForFirstScreen(stateName, yearOfConsult)
	if err != nil {
		log.Printf("[basic info state] error getting data for first screen(ano:%d, estado:%s):%q", yearOfConsult, stateName, err)
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetros ano=%d ou estado=%s são inválidos", yearOfConsult, stateName))
	}
	var agenciesBasic []models.AgencyBasic
	for k := range agencies {
		agenciesBasic = append(agenciesBasic, models.AgencyBasic{Name: agencies[k].ID, AgencyCategory: agencies[k].Entity})
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
	agencyMonthlyInfo, err := client.GetDataForSecondScreen(month, year, agencyName)
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

	members := map[int]int{10000: 0, 20000: 0, 30000: 0, 40000: 0, 50000: 0, -1: 0}
	servers := map[int]int{10000: 0, 20000: 0, 30000: 0, 40000: 0, 50000: 0, -1: 0}
	inactive := map[int]int{10000: 0, 20000: 0, 30000: 0, 40000: 0, 50000: 0, -1: 0}
	maxSalary := 0.0

	for _, employeeAux := range agencyMonthlyInfo.Employee {
		if employeeAux.Income.Total > maxSalary {
			maxSalary = employeeAux.Income.Total
		}
		salary := employeeAux.Income.Total
		var salaryRange int
		if salary <= 10000 {
			salaryRange = 10000
		} else if salary <= 20000 {
			salaryRange = 20000
		} else if salary <= 30000 {
			salaryRange = 30000
		} else if salary <= 40000 {
			salaryRange = 40000
		} else if salary <= 50000 {
			salaryRange = 50000
		} else {
			salaryRange = -1 // -1 is maker when the salary is over 50000
		}
		if employeeAux.Type == "membro" && employeeAux.Active == true {
			members[salaryRange]++
		} else if employeeAux.Type == "servidor" && employeeAux.Active == true {
			servers[salaryRange]++
		} else if employeeAux.Active == false {
			inactive[salaryRange]++
		}
	}

	return c.JSON(http.StatusOK, models.DataForChartAtAgencyScreen{
		Members:   members,
		Servers:   servers,
		Inactives: inactive,
		MaxSalary: maxSalary,
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
	agencyMonthlyInfo, err := client.GetDataForSecondScreen(month, year, agencyName)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d, mês=%d ou nome do orgão=%s são inválidos", year, month, agencyName))
	}

	fmt.Println()
	agencySummary := models.AgencySummary{
		TotalEmployees: agencyMonthlyInfo.Summary.General.Count,
		TotalWage:      agencyMonthlyInfo.Summary.General.Wage.Total,
		TotalPerks:     agencyMonthlyInfo.Summary.General.Perks.Total,
		MaxWage:        agencyMonthlyInfo.Summary.General.Wage.Max,
		CrawlingTime:   agencyMonthlyInfo.CrawlingTimestamp}
	return c.JSON(http.StatusOK, agencySummary)
}

var conf config

func main() {
	godotenv.Load() // There is no problem if the .env can not be loaded.

	err := envconfig.Process("remuneracao-magistrados", &conf)
	if err != nil {
		log.Fatal(err.Error())
	}

	dbClient, err := db.NewClient(conf.DBUrl, conf.DBName)
	if err != nil {
		log.Fatal(err)
	}
	defer dbClient.CloseConnection()

	// Criando o client do storage
	client, err = newClient(conf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Going to start listening at port:%d\n", conf.Port)

	e := echo.New()

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	// Enable access from all dadosjusbr domains.
	if os.Getenv("DADOSJUSBR_ENV") == "Prod" {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"https://dadosjusbr.com", "http://dadosjusbr.com", "https://dadosjusbr.org", "http://dadosjusbr.org"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderContentLength},
		}))
		log.Println("Using production CORS")
	} else {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderContentLength},
		}))
	}

	e.Renderer = renderer

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "ui/dist/",
		Browse: true,
		HTML5:  true,
		Index:  "index.html",
	}))

	e.Static("/static", "templates/assets")

	// Return a summary of an agency. This information will be used in the head of the agency page.
	e.GET("/uiapi/v1/orgao/resumo/:orgao/:ano/:mes", getSummaryOfAgency)
	// Return all the salary of a month and year. This will be used in the point chart at the entity page.
	e.GET("/uiapi/v1/orgao/salario/:orgao/:ano/:mes", getSalaryOfAgencyMonthYear)
	// Return the total of salary of every month of a year of a agency. The salary is divided in Wage, Perks and Others. This will be used to plot the bars chart at the state page.
	e.GET("/uiapi/v1/orgao/totais/:estado/:orgao/:ano", getTotalsOfAgencyYear)
	// Return basic information of a state
	e.GET("/uiapi/v1/orgao/:estado", getBasicInfoOfState)

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", conf.Port),
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	}
	e.Logger.Fatal(e.StartServer(s))
}
