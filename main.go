package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo"
)

type config struct {
	Port int `envconfig:"PORT"`
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

func main() {
	var conf config
	err := envconfig.Process("remuneracao-magistrados", &conf)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Going to start listening at port:%d\n", conf.Port)
	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", conf.Port),
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	}

	e := echo.New()

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("*.html")),
	}

	e.Renderer = renderer

	e.GET("/", handleDashboardRequest)
	e.Logger.Fatal(e.StartServer(s))
}
