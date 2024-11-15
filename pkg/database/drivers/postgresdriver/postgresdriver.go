package postgresdriver

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Raviraj2000/go-web-crawler/pkg/storage"
	_ "github.com/lib/pq"
)

type PostgresDriver struct {
	db *sql.DB
}

func (pg *PostgresDriver) Configure() error {
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")

	if host == "" || user == "" || password == "" || dbname == "" || port == "" {
		return fmt.Errorf("Postgres configuration is incomplete. Ensure all required environment variables are set")
	}

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to Postgres: %v", err)
	}

	pg.db = db
	return nil
}

func (pg *PostgresDriver) Save(data storage.PageData) error {
	query := "INSERT INTO pages (url, title, description) VALUES ($1, $2, $3)"
	_, err := pg.db.Exec(query, data.URL, data.Title, data.Desciption)
	return err
}

func (pg *PostgresDriver) Close() error {
	if pg.db != nil {
		return pg.db.Close()
	}
	return nil
}
