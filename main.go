package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dadosjusbr/api/papi"
	"github.com/dadosjusbr/api/uiapi"
	"github.com/dadosjusbr/storage"
	"github.com/dadosjusbr/storage/repo/database"
	"github.com/dadosjusbr/storage/repo/file_storage"
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

var pgS3Client *storage.Client
var loc *time.Location
var conf config

// newClient takes a config struct and creates a client to connect with DB and Cloud5
func newClient(db database.Interface, cloud file_storage.Interface) (*storage.Client, error) {
	client, err := storage.NewClient(db, cloud)
	if err != nil {
		return nil, fmt.Errorf("error creating storage.client: %q", err)
	}
	return client, nil
}

func newPostgresDB(c config) (*database.PostgresDB, error) {
	pgDb, err := database.NewPostgresDB(c.PgUser, c.PgPassword, c.PgDatabase, c.PgHost, c.PgPort)
	if err != nil {
		return nil, fmt.Errorf("error creating postgres DB client: %q", err)
	}
	return pgDb, nil
}

func newS3Client(c config) (*file_storage.S3Client, error) {
	s3Client, err := file_storage.NewS3Client(c.AwsRegion, c.AwsS3Bucket)
	if err != nil {
		return nil, fmt.Errorf("error creating S3 client: %v", err.Error())
	}
	return s3Client, nil
}

func main() {
	godotenv.Load() // There is no problem if the .env can not be loaded.
	l, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		log.Fatal(err.Error())
	}
	loc = l
	if err := envconfig.Process("", &conf); err != nil {
		log.Fatal(err.Error())
	}

	// Criando o client do S3
	s3Client, err := newS3Client(conf)
	if err != nil {
		log.Fatal(err)
	}

	pgDB, err := newPostgresDB(conf)
	if err != nil {
		log.Fatal(err)
	}
	pgS3Client, err = newClient(pgDB, s3Client)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := pgDB.GetConnection()
	if err != nil {
		log.Fatalf("Error connecting to postgres: %v", err)
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
	var nr *newrelic.Application
	if os.Getenv("DADOSJUSBR_ENV") == "Prod" {
		if conf.NewRelicApp == "" || conf.NewRelicLicense == "" {
			log.Fatalf("Missing environment variables NEWRELIC_APP_NAME or NEWRELIC_LICENSE")
		}
		nr, err = newrelic.NewApplication(
			newrelic.ConfigAppName(conf.NewRelicApp),
			newrelic.ConfigLicense(conf.NewRelicLicense),
			newrelic.ConfigAppLogForwardingEnabled(true),
		)
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
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderContentLength, echo.HeaderAccessControlAllowOrigin},
		}))
	}
	uiApiHandler, err := uiapi.NewHandler(*pgS3Client, conn, nr, conf.AwsRegion, conf.AwsS3Bucket, loc, conf.EnvOmittedFields, conf.SearchLimit, conf.DownloadLimit)
	if err != nil {
		log.Fatalf("Error creating uiapi handler: %q", err)
	}
	// Return a summary of an agency. This information will be used in the head of the agency page.
	uiAPIGroup.GET("/v1/orgao/resumo/:orgao/:ano/:mes", uiApiHandler.GetSummaryOfAgency)
	uiAPIGroup.GET("/v1/orgao/resumo/:orgao", uiApiHandler.GetAnnualSummary)
	// Return all the salary of a month and year. This will be used in the point chart at the entity page.
	uiAPIGroup.GET("/v1/orgao/salario/:orgao/:ano/:mes", uiApiHandler.GetSalaryOfAgencyMonthYear)
	// Return the total of salary of every month of a year of a agency. The salary is divided in Wage, Perks and Others. This will be used to plot the bars chart at the state page.
	uiAPIGroup.GET("/v1/orgao/totais/:orgao/:ano", uiApiHandler.GetTotalsOfAgencyYear)
	// Return basic information of a type or state
	uiAPIGroup.GET("/v1/orgao/:grupo", uiApiHandler.GetBasicInfoOfType)
	uiAPIGroup.GET("/v1/geral/remuneracao/:ano", uiApiHandler.GetGeneralRemunerationFromYear)
	uiAPIGroup.GET("/v1/geral/resumo", uiApiHandler.GeneralSummaryHandler)
	// Retorna um conjunto de dados a partir de filtros informados por query params
	uiAPIGroup.GET("/v2/pesquisar", uiApiHandler.SearchByUrl)
	// Baixa um conjunto de dados a partir de filtros informados por query params
	uiAPIGroup.GET("/v2/download", uiApiHandler.DownloadByUrl)

	apiHandler := papi.NewHandler(*pgS3Client, conf.DadosJusURL, conf.PackageRepoURL)
	// Public API configuration
	apiGroup := e.Group("/v1", middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderContentLength},
	}))
	// Return agency
	apiGroup.GET("/orgao/:orgao", apiHandler.GetAgencyById)
	// Return all agencies
	apiGroup.GET("/orgaos", apiHandler.GetAllAgencies)
	// Return MIs by year
	apiGroup.GET("/dados/:orgao/:ano", apiHandler.GetMonthlyInfo)
	// Return MIs by month
	apiGroup.GET("/dados/:orgao/:ano/:mes", apiHandler.GetMonthlyInfo)

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", conf.Port),
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	}
	e.Logger.Fatal(e.StartServer(s))
}
