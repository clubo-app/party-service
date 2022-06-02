package repository

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	g "github.com/golang-migrate/migrate/v4/source/github"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewPGXPool(dbUser, dbPW, dbName, dbHost string, dbPort uint16) (*pgxpool.Pool, error) {
	urlStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPW, dbHost, fmt.Sprint(dbPort), dbName)
	pgURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	connURL := *pgURL
	if connURL.Scheme == "cockroachdb" {
		connURL.Scheme = "postgres"
	}

	c, err := pgxpool.ParseConfig(connURL.String())
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), c)
	if err != nil {
		return nil, fmt.Errorf("pgx connection error: %w", err)
	}

	err = validateSchema(connURL)
	if err != nil {
		log.Printf("Schema validation error: %v", err)
	}

	return pool, nil
}

const version = 1

func validateSchema(url url.URL) error {
	url.Scheme = "pgx"
	url2 := fmt.Sprintf("%v%v", url.String(), "?sslmode=disable")
	g := g.Github{}
	d, err := g.Open("github://clubo-app/party-service/repository/migrations")
	if err != nil {
		return err
	}
	defer d.Close()

	m, err := migrate.NewWithSourceInstance("github", d, url2)

	if err != nil {
		return err
	}
	err = m.Migrate(version) // current version
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	defer m.Close()
	return nil
}
