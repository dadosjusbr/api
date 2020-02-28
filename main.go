package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dadosjusbr/remuneracao-magistrados/db"
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
	MongoMICol  string `envconfig:"MONGODB_MICOL"`
	MongoAgCol  string `envconfig:"MONGODB_AGCOL"`
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
	var monthTotalsOfYear []monthTotals
	agencyName := c.Param("orgao")
	for _, agencyMonthlyInfo := range agenciesMonthlyInfo[agencyName] {
		monthTotals := monthTotals{agencyMonthlyInfo.Month, agencyMonthlyInfo.Summary.Wage.Total, agencyMonthlyInfo.Summary.Perks.Total, agencyMonthlyInfo.Summary.Others.Total}
		monthTotalsOfYear = append(monthTotalsOfYear, monthTotals)
	}
	agencyTotalsYear := agencyTotalsYear{year, monthTotalsOfYear}
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
	var agenciesBasic []agencyBasic
	for k := range agencies {
		agenciesBasic = append(agenciesBasic, agencyBasic{agencies[k].ID, agencies[k].Entity})
	}
	state := state{stateName, "", "", agenciesBasic}
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
	var employees []employee
	for _, employeeAux := range agencyMonthlyInfo.Employee {
		newEmployee := employee{
			employeeAux.Name,
			*employeeAux.Income.Wage,
			employeeAux.Income.Perks.Total,
			employeeAux.Income.Other.Total,
			employeeAux.Income.Total}
		employees = append(employees, newEmployee)
	}
	return c.JSON(http.StatusOK, employees)
}

func getSummaryOfAgency(c echo.Context) error {
	yearOfCosult := time.Now().Year()
	monthOfConsult := 1 //int(time.Now().Month()) Tem um erro na api de leitura enquanto não ajeitar deixei hardcoded aqui.
	agencyName := c.Param("orgao")
	var agencyMonthlyInfo *storage.AgencyMonthlyInfo
	var err error
	for i := monthOfConsult; i > 0; i-- {
		agencyMonthlyInfo, err = client.GetDataForSecondScreen(monthOfConsult, yearOfCosult, agencyName)
		if err == nil {
			break
		}
	}
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Parâmetro ano=%d, mês=%d ou nome do orgão=%s são inválidos", yearOfCosult, monthOfConsult, agencyName))
	}
	agencySummary := agencySummary{agencyMonthlyInfo.Summary.Count, agencyMonthlyInfo.Summary.Wage.Total, agencyMonthlyInfo.Summary.Perks.Total, agencyMonthlyInfo.Summary.Wage.Max}
	return c.JSON(http.StatusOK, agencySummary)
}

func main() {
	godotenv.Load() // There is no problem if the .env can not be loaded.
	var conf config
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

	e.Use(middleware.CORS())

	e.Renderer = renderer

	e.Static("/static", "templates/assets")
	e.Static("/novo", "ui/dist")

	e.GET("/", handleMainPageRequest(dbClient))
	e.GET("/:year/:month", handleMonthRequest(dbClient))

	// Return a summary of an agency. This information will be used in the head of the agency page.
	e.GET("/uiapi/v1/orgao/resumo/:orgao", getSummaryOfAgency)
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

type state struct {
	Name      string
	ShortName string
	FlagURL   string
	Agency    []agencyBasic
}

type agencyBasic struct {
	Name           string
	AgencyCategory string
}

type employee struct {
	Name   string
	Wage   float64
	Perks  float64
	Others float64
	Total  float64
}

type agencySummary struct {
	TotalEmployees int
	TotalWage      float64
	TotalPerks     float64
	MaxWage        float64
}

type agencyTotalsYear struct {
	Year        int
	MonthTotals []monthTotals
}

type monthTotals struct {
	Month  int
	Wage   float64
	Perks  float64
	Others float64
}
