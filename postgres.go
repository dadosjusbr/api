package main

import (
	"context"
	"fmt"
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
