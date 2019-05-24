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
)

type config struct {
	Port  int    `envconfig:"PORT"`
	DBUrl string `envconfig:"MONGODB_URI"`
}

// Entry represents each entry of the processed data.
type Entry struct {
	PCloudURL   string
	DataHubURL  string
	Success     bool
	Indentifier string
}

// EntryViewModel is the one that the template is going to use.
type EntryViewModel struct {
	Item            Entry
	PreviousEntries []Entry
	Month           int
	Year            int
}

func loadEntryByMonthAndYear(month int, year int) (Entry, error) {
	entry := Entry{
		PCloudURL:   "https://my.pcloud.com/publink/show?code=XZ7fM17Z6S7V93BsYXhMaT4irtMMO8kXc1IV",
		DataHubURL:  "https://my.pcloud.com/publink/show?code=XZ7fM17Z6S7V93BsYXhMaT4irtMMO8kXc1IV",
		Success:     true,
		Indentifier: strconv.Itoa(month) + "-" + strconv.Itoa(year),
	}

	return entry, nil
}

var monthsLabelMap = map[int]string{
	1:  "Janeiro",
	2:  "Fevereiro",
	3:  "Marco",
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

func loadPreviousEntries() ([2]Entry, error) {
	var entries [2]Entry

	entryA, _ := loadEntryByMonthAndYear(3, 2018)
	entryB, _ := loadEntryByMonthAndYear(2, 2018)

	entries[0] = entryA
	entries[1] = entryB

	return entries, nil
}

func handleDashboardRequest(c echo.Context) error {
	month, year := 4, 2018
	entry, _ := loadEntryByMonthAndYear(month, year)
	oldEntries, _ := loadPreviousEntries()

	data := EntryViewModel{
		Item:            entry,
		Month:           month,
		Year:            year,
		PreviousEntries: oldEntries[:],
	}

	return c.Render(http.StatusOK, "index.html", data)
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

func getHandleMonthRequest(dbClient *db.Client) echo.HandlerFunc {
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
			return c.String(http.StatusInternalServerError, "unexpected error")
		}

		monthLabel := fmt.Sprintf("%s %d", monthsLabelMap[month], year)

		processedMonths, err := dbClient.GetProcessedMonths()
		if err != nil {
			fmt.Println(fmt.Errorf("error retrieving all parsed months from db --> %v", err))
			return c.String(http.StatusInternalServerError, "unexpected error")
		}

		var sidebarElements []SidebarElement

		for _, pm := range processedMonths {
			label := fmt.Sprintf("%s %d", monthsLabelMap[pm.Month], pm.Year)
			URL := fmt.Sprintf("/%d/%d", pm.Year, pm.Month)
			sidebarElements = append(sidebarElements, SidebarElement{Label: label, URL: URL})
		}

		viewModel := struct {
			Month           int
			Year            int
			MonthLabel      string
			SpreadsheetsURL string
			DatapackageURL  string
			SidebarElements []SidebarElement
		}{
			monthResults.Month,
			monthResults.Year,
			monthLabel,
			monthResults.SpreadsheetsURL,
			monthResults.DatapackageURL,
			sidebarElements,
		}
		return c.Render(http.StatusOK, "monthTemplate.html", viewModel)
	}
}

func main() {
	var conf config
	err := envconfig.Process("remuneracao-magistrados", &conf)
	if err != nil {
		log.Fatal(err.Error())
	}

	dbClient, err := db.NewClient(conf.DBUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer dbClient.CloseConnection()

	fmt.Printf("Going to start listening at port:%d\n", conf.Port)

	e := echo.New()

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	e.Renderer = renderer

	e.Static("/static", "templates/assets")

	e.GET("/", handleDashboardRequest)
	e.GET("/:year/:month", getHandleMonthRequest(dbClient))

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", conf.Port),
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	}
	e.Logger.Fatal(e.StartServer(s))
}
