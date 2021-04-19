package postgres

import (
	"database/sql"
	"embed"
	"fmt"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/trelore/todoapi/internal"

	// File library for reading migrations
	_ "github.com/golang-migrate/migrate/v4/source/file"
	// Postgres library
	_ "github.com/lib/pq"
)

type p struct {
	db *sql.DB
}

//go:embed migrations/*.sql
var migrationsFS embed.FS

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME     = "postgres"
)

func New() (*p, func() error, error) {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return nil, nil, fmt.Errorf("opening postgres: %w", err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{MigrationsTable: "migrations"})
	if err != nil {
		return nil, nil, fmt.Errorf("with instance: %w", err)
	}
	source, err := httpfs.New(http.FS(migrationsFS), "migrations")
	if err != nil {
		return nil, nil, fmt.Errorf("getting migration files: %w", err)
	}
	mig, err := migrate.NewWithInstance("httpfs", source, DB_NAME, driver)
	if err != nil {
		return nil, nil, fmt.Errorf("migration instance: %w", err)
	}
	err = mig.Up()
	if err != nil {
		return nil, nil, fmt.Errorf("migrating up: %w", err)
	}
	return &p{
		db: db,
	}, db.Close, nil
}

// Delete implements the interface
func (p *p) Delete(id string) error {
	return nil
}

// Insert implements the interface
func (p *p) Insert(description string) (*internal.Item, error) {
	return nil, nil
}

// List implements the interface
func (p *p) List() ([]*internal.Item, error) {
	return nil, nil
}

// Get implements the interface
func (p *p) Get(id string) (*internal.Item, error) {
	return nil, nil
}

// Upsert implements the interface
func (p *p) Upsert(id string, item *internal.Item) (*internal.Item, error) {
	return nil, nil
}
