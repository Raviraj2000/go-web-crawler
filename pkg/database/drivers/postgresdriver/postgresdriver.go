package postgresdriver

import (
	"database/sql"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

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

	err := startPortForwarding("postgres-service", "5432", "5432", "default") // Change "default" to your namespace if needed
	if err != nil {
		return nil, fmt.Errorf("failed to start port-forwarding: %v", err)
	}

	// Step 2: Wait for port-forwarding to establish
	time.Sleep(5 * time.Second) // Adjust based on your environment

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.User, config.Password, config.DBName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &PostgresDriver{db: db}, nil
}

// Save inserts PageData into the Postgres database.
func (p *PostgresDriver) Save(data models.PageData) error {
	_, err := p.db.Exec("INSERT INTO pages (url, title, description) VALUES ($1, $2, $3)",
		data.URL, data.Title, data.Description)
	return err
}

func startPortForwarding(service, localPort, remotePort, namespace string) error {
	// Prepare the kubectl port-forward command
	cmd := exec.Command("kubectl", "port-forward", fmt.Sprintf("svc/%s", service), fmt.Sprintf("%s:%s", localPort, remotePort), "-n", namespace)

	// Run the command in the background
	cmd.Stdout = nil
	cmd.Stderr = nil

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start kubectl port-forward: %v", err)
	}

	go func() {
		if err := cmd.Wait(); err != nil {
			if !strings.Contains(err.Error(), "signal: killed") {
				log.Printf("port-forward command exited: %v", err)
			}
		}
	}()

	log.Printf("Port-forwarding to service %s started on port %s", service, localPort)
	return nil
}

// Close closes the database connection.
func (p *PostgresDriver) Close() error {
	return p.db.Close()
}
