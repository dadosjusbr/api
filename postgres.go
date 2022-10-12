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

func (p PostgresDB) Filter(query string, arguments []interface{}) ([]models.SearchResult, error) {
	results := []models.SearchResult{}
	var err error
	txn := p.newrelic.StartTransaction("pg.Filter")
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

func (p PostgresDB) LowCostFilter(query string, arguments []interface{}) ([]models.SearchDetails, error) {
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

func (p PostgresDB) Count(query string, arguments []interface{}) (int, error) {
	var count int
	var err error

	txn := p.newrelic.StartTransaction("pg.Count")
	defer txn.End()
	ctx := newrelic.NewContext(context.Background(), txn)
	if len(arguments) > 0 {
		err = p.conn.QueryRowContext(ctx, query, arguments...).Scan(&count)
	} else {
		err = p.conn.QueryRowContext(ctx, query).Scan(&count)
	}
	if err != nil {
		return -1, fmt.Errorf("erro ao pegar contagem de resultados: %v", err)
	}
	return count, nil
}

func (p PostgresDB) RemunerationQuery(filter *models.Filter, limit int) string {
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
		p.AddFiltersInQuery(&query, filter)
	}

	return fmt.Sprintf("%s FETCH FIRST %d ROWS ONLY;", query, limit)
}

//Função que recebe os filtros e a partir deles estrutura a query SQL da pesquisa
func (p PostgresDB) LowCostRemunerationQuery(filter *models.Filter) string {
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
		p.LowCostAddFiltersInQuery(&query, filter)
	}

	return query
}

//Função que insere os filtros na query
func (p PostgresDB) LowCostAddFiltersInQuery(query *string, filter *models.Filter) {
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

func (p PostgresDB) AddFiltersInQuery(query *string, filter *models.Filter) {
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
	if filter.Category != "" {
		lastIndex := len(filter.Years) + len(filter.Months) + len(filter.Agencies)
		if lastIndex > 0 {
			*query = fmt.Sprintf("%s AND", *query)
		}
		*query = fmt.Sprintf("%s r.categoria_contracheque = $%d", *query, lastIndex+1)
	}
}

//Função que define os argumentos passados para a query
func (p PostgresDB) LowCostArguments(filter *models.Filter) []interface{} {
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
		if filter.Category != "" {
			arguments = append(arguments, filter.Category)
		}
		if filter.Types != "" {
			// Adicionando '% %' na clausura LIKE
			arguments = append(arguments, fmt.Sprintf("%%%s%%", filter.Types))
		}
	}

	return arguments
}
