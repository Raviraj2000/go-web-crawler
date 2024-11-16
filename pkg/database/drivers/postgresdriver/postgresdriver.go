package postgresdriver

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Raviraj2000/go-web-crawler/pkg/storage/models"
	_ "github.com/lib/pq"
)

type PostgresDriver struct {
	db *sql.DB
}

// Config represents the configuration required to initialize a PostgresDriver instance.
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewPostgresDriver initializes a new PostgresDriver instance.
func NewPostgresDriver(config Config) (models.StorageDriver, error) {

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)
	log.Printf("Connecting to database: %s", dsn)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}
	// Validate connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	err = createPagesTable(db)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	log.Println("Table 'pages' is ready.")

	return &PostgresDriver{db: db}, nil
}

func createPagesTable(db *sql.DB) error {
	// SQL to create the table
	query := `
	CREATE TABLE IF NOT EXISTS pages (
		id SERIAL PRIMARY KEY,
		url TEXT NOT NULL,
		title TEXT,
		description TEXT
	);`

	// Execute the query
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}

	return nil
}

// Save inserts PageData into the Postgres database.
func (p *PostgresDriver) Save(data models.PageData) error {
	_, err := p.db.Exec("INSERT INTO pages (url, title, description) VALUES ($1, $2, $3)",
		data.URL, data.Title, data.Description)
	return err
}

// Close closes the database connection.
func (p *PostgresDriver) Close() error {
	return p.db.Close()
}
