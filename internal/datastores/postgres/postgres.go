package postgres

import (
	"database/sql"
	"embed"
	"errors"
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
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)

	db, err := sql.Open("postgres", dsn)
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
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, nil, fmt.Errorf("migrating up: %w", err)
	}
	return &p{
		db: db,
	}, db.Close, nil
}

// Delete implements the interface
func (p *p) Delete(id string) error {
	query := `DELETE FROM todo WHERE id = $1`
	_, err := p.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("deleting item: %w", err)
	}
	return nil
}

// Insert implements the interface
func (p *p) Insert(description string) (*internal.Item, error) {
	item := internal.Item{}
	query := `INSERT INTO todo (description) VALUES ($1) RETURNING id, done`
	err := p.db.QueryRow(query, description).Scan(&item.ID, &item.Done)
	if err != nil {
		return nil, fmt.Errorf("inserting item: %w", err)
	}
	item.Description = description
	return &item, nil
}

// List implements the interface
func (p *p) List() ([]*internal.Item, error) {
	query := `SELECT id, description, done FROM todo`
	rows, err := p.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("querying items: %w", err)
	}
	defer rows.Close()
	var items []*internal.Item
	for rows.Next() {
		var item internal.Item
		if err = rows.Scan(&item.ID, &item.Description, &item.Done); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating result set: %w", err)
	}

	return items, nil
}

// Get implements the interface
func (p *p) Get(id string) (*internal.Item, error) {
	var item internal.Item
	query := `SELECT id, description, done FROM todo WHERE id = $1`
	err := p.db.QueryRow(query, id).Scan(&item.ID, &item.Description, &item.Done)
	if err != nil {
		return nil, fmt.Errorf("querying row: %w", err)
	}
	return &item, nil
}

// Upsert implements the interface
func (p *p) Upsert(id string, item *internal.Item) (*internal.Item, error) {
	const query = `UPDATE todo SET description = $2, done=$3 WHERE id = $1 RETURNING id, description, done`
	err := p.db.QueryRow(query, id, item.Description, item.Done).Scan(&item.ID, &item.Description, &item.Done)
	if err != nil {
		return nil, fmt.Errorf("upserting row: %w", err)
	}
	return item, nil
}
