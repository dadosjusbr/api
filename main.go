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

	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type config struct {
	Port   int    `envconfig:"PORT"`
	DBUrl  string `envconfig:"MONGODB_URI"`
	DBName string `envconfig:"MONGODB_NAME"`
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
			fmt.Println(err)
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

func getTotalsOfAgencyYear(c echo.Context) error {
	monthTotals1 := monthTotals{1, 100000.0, 25000.0, 65000.0}
	monthTotals2 := monthTotals{2, 150000.0, 35000.0, 55000.0}
	monthTotals3 := monthTotals{3, 120000.0, 28000.0, 49000.0}
	agencyTotalsYear := agencyTotalsYear{2018, []monthTotals{monthTotals1, monthTotals2, monthTotals3}}
	return c.JSON(http.StatusOK, agencyTotalsYear)
}

func getSummaryOfEntitiesOfState(c echo.Context) error {
	agencyBasic1 := agencyBasic{"TJPB", "J"}
	agencyBasic2 := agencyBasic{"MPPB", "M"}
	agencyBasic3 := agencyBasic{"TRTPB", "J"}
	state := state{"Paraíba", "pb", "url", []agencyBasic{agencyBasic1, agencyBasic2, agencyBasic3}}
	return c.JSON(http.StatusOK, state)
}

func getSalaryOfAgencyMonthYear(c echo.Context) error {
	employee1 := employee{"Marcos", 30000.0, 14000.0, 25000.0, 69000.0}
	employee2 := employee{"Joeberth", 35000.0, 19000.0, 20000.0, 74000.0}
	employee3 := employee{"Maria", 34000.0, 15000.0, 23000.0, 72000.0}
	employees := []employee{employee1, employee2, employee3}
	return c.JSON(http.StatusOK, employees)
}

func getSummaryOfAgency(c echo.Context) error {
	agencySummary := agencySummary{100, 250000.0, 100000.0, 26000.0}
	return c.JSON(http.StatusOK, agencySummary)
}

func main() {
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
	// Return all the salary of a month and year. This will be used in the point chart at the agency page.
	e.GET("/uiapi/v1/orgao/salario/:orgao/:year/:month", getSalaryOfAgencyMonthYear)
	// This will return information of a state and its entities and agencies. This will be used to provide basic information for the state page.
	e.GET("/uiapi/v1/entidades/resumo/:estado", getSummaryOfEntitiesOfState)
	// Return the total of salary of every month of a year of a agency. The salary is divided in Wage, Perks and Others. This will be used to plot the bars chart at the state page.
	e.GET("/uiapi/v1/orgao/totais/:orgao/:year", getTotalsOfAgencyYear)

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

type agency struct {
	Name           string
	ShortName      string
	AgencyCategory string
	AgencySummary  agencySummary
	Employee       []employee
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
