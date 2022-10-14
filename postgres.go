package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dadosjusbr/api/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/newrelic/go-agent/v3/integrations/nrpq"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type PostgresDB struct {
	conn     *sqlx.DB
	newrelic *newrelic.Application
}

type PostgresCredentials struct {
	user     string
	password string
	dbName   string
	host     string
	port     string
	uri      string
}

//Retorna uma nova conexão com o postgres, através da uri passada como parâmetro
func NewPostgresDB(pgCredentials PostgresCredentials) (*PostgresDB, error) {
	conn, err := sqlx.Open("nrpostgres", pgCredentials.uri)
	if err != nil {
		return nil, fmt.Errorf("error while accessing database: %q", err)
	}
	ctx, canc := context.WithTimeout(context.Background(), 30*time.Second)
	defer canc()
	if err := conn.DB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("Error connecting to postgres (creds:%+v):%w", pgCredentials, err)
	}
	return &PostgresDB{
		conn: conn,
	}, nil
}

//Recebe os dados da conexão, verifica se está tudo certo e depois retorna a uri da conexão
func NewPgCredentials(c config) (*PostgresCredentials, error) {
	for k, v := range map[string]string{
		"postgres-user":     c.PgUser,
		"postgres-password": c.PgPassword,
		"postgres-database": c.PgDatabase,
		"postgres-host":     c.PgHost,
		"postgres-port":     c.PgPort,
	} {
		if v == "" {
			return nil, fmt.Errorf("%s is not set!", k)
		}
	}
	return &PostgresCredentials{
		c.PgUser,
		c.PgPassword,
		c.PgDatabase,
		c.PgHost,
		c.PgPort,
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			c.PgUser,
			c.PgPassword,
			c.PgHost,
			c.PgPort,
			c.PgDatabase),
	}, nil
}

func (p *PostgresDB) Disconnect() error {
	err := p.conn.Close()
	if err != nil {
		return fmt.Errorf("error closing connection: %q", err)
	}
	return nil
}

func (p PostgresDB) Filter(query string, arguments []interface{}) ([]models.SearchDetails, error) {
	results := []models.SearchDetails{}
	var err error
	txn := p.newrelic.StartTransaction("pg.LowCostFilter")
	defer txn.End()
	ctx := newrelic.NewContext(context.Background(), txn)
	if len(arguments) > 0 {
		err = p.conn.SelectContext(ctx, &results, query, arguments...)
	} else {
		err = p.conn.SelectContext(ctx, &results, query)
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer a seleção por filtro: %v", err)
	}
	return results, nil
}

//Função que recebe os filtros e a partir deles estrutura a query SQL da pesquisa
func (p PostgresDB) RemunerationQuery(filter *models.Filter) string {
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
	if filter != nil {
		p.AddFiltersInQuery(&query, filter)
	}

	return query
}

//Função que insere os filtros na query
func (p PostgresDB) AddFiltersInQuery(query *string, filter *models.Filter) {
	*query = *query + " WHERE"

	//Insere os filtros de ano caso existam
	if len(filter.Years) > 0 {
		var years []string
		years = append(years, filter.Years...)
		for i := 0; i < len(filter.Years); i++ {
			years[i] = fmt.Sprintf("$%d", i+1)
		}
		*query = fmt.Sprintf("%s ano IN (%s)", *query, strings.Join(years, ","))
	}

	//Insere os filtros de mês
	if len(filter.Months) > 0 {
		lastIndex := len(filter.Years)
		if lastIndex > 0 {
			*query = fmt.Sprintf("%s AND", *query)
		}
		var months []string
		months = append(months, filter.Months...)
		for i := lastIndex; i < len(filter.Months)+lastIndex; i++ {
			months[i-lastIndex] = fmt.Sprintf("$%d", i+1)
		}
		*query = fmt.Sprintf("%s mes IN (%s)", *query, strings.Join(months, ","))
	}

	//Insere o filtro de órgãos
	if len(filter.Agencies) > 0 {
		lastIndex := len(filter.Years) + len(filter.Months)
		if lastIndex > 0 {
			*query = fmt.Sprintf("%s AND", *query)
		}
		var agencies []string
		agencies = append(agencies, filter.Agencies...)
		for i := lastIndex; i < lastIndex+len(filter.Agencies); i++ {
			agencies[i-lastIndex] = fmt.Sprintf("$%d", i+1)
		}
		*query = fmt.Sprintf("%s id_orgao IN (%s)", *query, strings.Join(agencies, ","))
	}
}

//Função que define os argumentos passados para a query
func (p PostgresDB) Arguments(filter *models.Filter) []interface{} {
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
	}

	return arguments
}

//Contando a quantidade de órgãos que temos no banco
func (p PostgresDB) CountAgencies() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM orgaos`
	err := p.conn.Get(&count, query)
	if err != nil {
		return 0, fmt.Errorf("error counting agencies: %q", err)
	}
	return count, nil
}

//Contando a quantidade de coletas(registros de remunerações) que temos no banco
func (p PostgresDB) CountRemunerationRecords() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM coletas`
	err := p.conn.Get(&count, query)
	if err != nil {
		return 0, fmt.Errorf("error counting collections: %q", err)
	}
	return count, nil
}

/*Pegando o primeiro registro de remuneração que temos.
A query ordena os registros por ano e mês e pega o primeiro registro*/
func (p PostgresDB) GetFirstDateWithRemunerationRecords() (time.Time, error){
	query := `SELECT ano, mes FROM coletas ORDER BY ano, mes LIMIT 1`
	var year, month int
	err := p.conn.QueryRow(query).Scan(&year, &month)
	if err != nil {
		return time.Time{}, fmt.Errorf("error getting first date with remuneration records: %q", err)
	}
	return time.Date(year, time.Month(month), 2, 0, 0, 0, 0, time.UTC).In(loc), nil
}

/*Pegando o último registro de remuneração que temos.
A query ordena, em ordem decrescente, os registros por ano e mês e pega 
o primeiro registro*/
func (p PostgresDB) GetLastDateWithRemunerationRecords() (time.Time, error){
	query := `SELECT ano, mes FROM coletas ORDER BY ano DESC, mes DESC LIMIT 1`
	var year, month int
	err := p.conn.QueryRow(query).Scan(&year, &month)
	if err != nil {
		return time.Time{}, fmt.Errorf("error getting last date with remuneration records: %q", err)
	}
	return time.Date(year, time.Month(month), 2, 0, 0, 0, 0, time.UTC).In(loc), nil
}

/*Pegando a soma de todas as remunerações*/
func (p PostgresDB) GetGeneralRemunerationValue() (float64, error){
	var value float64
	query := `SELECT SUM(CAST(sumario -> 'base_remuneration' ->> 'total' AS DECIMAL) + CAST(sumario -> 'other_remunerations' ->> 'total' AS DECIMAL)) FROM coletas;`
	err := p.conn.Get(&value, query)
	if err != nil {
		return 0, fmt.Errorf("error getting general remuneration value: %q", err)
	}
	return value, nil
}