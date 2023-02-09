package uiapi

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/newrelic/go-agent/v3/integrations/nrpq"
	"github.com/newrelic/go-agent/v3/newrelic"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresDB struct {
	conn        *gorm.DB
	newrelic    *newrelic.Application
	credentials PostgresCredentials
}

type PostgresCredentials struct {
	user     string
	password string
	dbName   string
	host     string
	port     string
	uri      string
}

// Recebe os dados da conexão, verifica se está tudo certo e depois retorna a uri da conexão
func NewPgCredentials(user, password, dbName, host, port string) (*PostgresCredentials, error) {
	for k, v := range map[string]string{
		"postgres-user":     user,
		"postgres-password": password,
		"postgres-database": dbName,
		"postgres-host":     host,
		"postgres-port":     port,
	} {
		if v == "" {
			return nil, fmt.Errorf("%s is not set!", k)
		}
	}
	return &PostgresCredentials{
		user,
		password,
		dbName,
		host,
		port,
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			user,
			password,
			host,
			port,
			dbName),
	}, nil
}

// Retorna uma nova conexão com o postgres, através da uri passada como parâmetro
func newPostgresDB(pgCredentials PostgresCredentials) (*postgresDB, error) {
	conn, err := sql.Open("nrpostgres", pgCredentials.uri)
	if err != nil {
		panic(err)
	}
	ctx, canc := context.WithTimeout(context.Background(), 30*time.Second)
	defer canc()
	if err := conn.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error connecting to postgres (creds:%s):%q", pgCredentials.uri, err)
	}
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: conn,
	}))
	if err != nil {
		return nil, fmt.Errorf("error initializing gorm: %q", err)
	}
	return &postgresDB{
		conn: db,
	}, nil
}

func (p *postgresDB) connect() error {
	if p.conn != nil {
		return nil
	} else {
		conn, err := sql.Open("nrpostgres", p.credentials.uri)
		if err != nil {
			panic(err)
		}
		ctx, canc := context.WithTimeout(context.Background(), 30*time.Second)
		defer canc()
		if err := conn.PingContext(ctx); err != nil {
			return fmt.Errorf("error connecting to postgres (creds:%s):%q", p.credentials.uri, err)
		}
		db, err := gorm.Open(postgres.New(postgres.Config{
			Conn: conn,
		}))
		if err != nil {
			return fmt.Errorf("error initializing gorm: %q", err)
		}
		p.conn = db
		return nil
	}
}

func (p *postgresDB) disconnect() error {
	db, err := p.conn.DB()
	if err != nil {
		return fmt.Errorf("error returning sql DB: %q", err)
	}
	err = db.Close()
	if err != nil {
		return fmt.Errorf("error closing DB connection: %q", err)
	}
	return nil
}

func (p postgresDB) filter(query string, arguments []interface{}) ([]searchDetails, error) {
	results := []searchDetails{}
	var err error
	txn := p.newrelic.StartTransaction("pg.LowCostFilter")
	defer txn.End()
	ctx := newrelic.NewContext(context.Background(), txn)
	if len(arguments) > 0 {
		err = p.conn.WithContext(ctx).Raw(query, arguments...).Scan(&results).Error
	} else {
		err = p.conn.WithContext(ctx).Raw(query).Scan(&results).Error
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer a seleção por filtro: %v", err)
	}
	return results, nil
}

// Função que recebe os filtros e a partir deles estrutura a query SQL da pesquisa
func (p postgresDB) remunerationQuery(searchParams *searchParams) string {
	//A query padrão sem os filtros
	query := `SELECT
		id_orgao as orgao,
		mes as mes,
		ano as ano,
		linhas_descontos as descontos,
		linhas_base as base,
		linhas_outras as outras,
		zip_url as zip_url
	FROM remuneracoes_zips 
	`
	if searchParams != nil {
		p.addFiltersInQuery(&query, searchParams)
	}

	return query
}

// Função que insere os filtros na query
func (p postgresDB) addFiltersInQuery(query *string, searchParams *searchParams) {
	*query = *query + " WHERE"

	//Insere os filtros de ano caso existam
	if len(searchParams.Years) > 0 {
		var years []string
		years = append(years, searchParams.Years...)
		for i := 0; i < len(searchParams.Years); i++ {
			years[i] = fmt.Sprintf("$%d", i+1)
		}
		*query = fmt.Sprintf("%s ano IN (%s)", *query, strings.Join(years, ","))
	}

	//Insere os filtros de mês
	if len(searchParams.Months) > 0 {
		lastIndex := len(searchParams.Years)
		if lastIndex > 0 {
			*query = fmt.Sprintf("%s AND", *query)
		}
		var months []string
		months = append(months, searchParams.Months...)
		for i := lastIndex; i < len(searchParams.Months)+lastIndex; i++ {
			months[i-lastIndex] = fmt.Sprintf("$%d", i+1)
		}
		*query = fmt.Sprintf("%s mes IN (%s)", *query, strings.Join(months, ","))
	}

	//Insere o filtro de órgãos
	if len(searchParams.Agencies) > 0 {
		lastIndex := len(searchParams.Years) + len(searchParams.Months)
		if lastIndex > 0 {
			*query = fmt.Sprintf("%s AND", *query)
		}
		var agencies []string
		agencies = append(agencies, searchParams.Agencies...)
		for i := lastIndex; i < lastIndex+len(searchParams.Agencies); i++ {
			agencies[i-lastIndex] = fmt.Sprintf("$%d", i+1)
		}
		*query = fmt.Sprintf("%s id_orgao IN (%s)", *query, strings.Join(agencies, ","))
	}
}

// Função que define os argumentos passados para a query
func (p postgresDB) arguments(searchParams *searchParams) []interface{} {
	var arguments []interface{}
	if searchParams != nil {
		if len(searchParams.Years) > 0 {
			for _, y := range searchParams.Years {
				arguments = append(arguments, y)
			}
		}
		if len(searchParams.Months) > 0 {
			for _, m := range searchParams.Months {
				arguments = append(arguments, m)
			}
		}
		if len(searchParams.Agencies) > 0 {
			for _, a := range searchParams.Agencies {
				arguments = append(arguments, a)
			}
		}
	}

	return arguments
}
