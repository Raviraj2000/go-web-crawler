package storage

import (
	"errors"

	"github.com/Raviraj2000/go-web-crawler/pkg/database/drivers/postgresdriver"
	"github.com/Raviraj2000/go-web-crawler/pkg/storage/models"
)

// DriverFactory initializes the appropriate storage driver based on driverType.
func DriverFactory(driverType string) (models.StorageDriver, error) {
	switch driverType {
	case "postgres":
		config := postgresdriver.Config{
			Host:     "localhost",
			Port:     5432,
			User:     "root",
			Password: "secret",
			DBName:   "WebCrawlerDB",
			SSLMode:  "disable",
		}
		return postgresdriver.NewPostgresDriver(config)
	default:
		return nil, errors.New("unsupported storage driver type")
	}
}
