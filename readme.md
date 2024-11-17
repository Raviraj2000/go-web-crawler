# Web Crawler with Dynamic Storage Driver Support

A modular and scalable web crawler built with Go, designed to support multiple storage backends dynamically. The project is deployable using Kubernetes and is powered by Redis for queue management and URL deduplication.

---

## Features

- **Dynamic Storage Backend:** Easily switch between storage drivers like PostgreSQL or implement your own custom driver.
- **Scalable Architecture:** Kubernetes deployment for distributed crawling with multiple pods.
- **URL Deduplication:** Redis-based deduplication to avoid re-crawling URLs.
- **Highly Configurable:** Control crawling behavior via environment variables or Kubernetes ConfigMaps.
- **Rate Limiting:** Rate-limited requests to avoid overwhelming servers.

## Prerequisites

- Docker for containerization.
- Kubernetes (e.g., Minikube or any cluster) for orchestration.
- Redis for URL queue and deduplication.

## Setup and Deployment

1. Clone the Repository
    ```sh
    git clone https://github.com/Raviraj2000/go-web-crawler.git
    cd go-web-crawler
    ```

2. Use the provided `deploy` script for deployment:
   
   For Windows:
    ```sh
    ./deploy.ps1 -WaitTime 60 -build $true -SeedUrl "https://example.com" -WorkerCount "100" -MaxUrls "10000" -StorageDriver "postgres"
    ```

    For Mac and linux:

    ```sh
    ./deploy.sh --wait-time 60 --build --seed-url "https://example.com" --worker-count 100 --max-urls 10000 --storage-driver postgres
    ```

    | Parameter      | Description                          | Example Value         |
    |----------------|--------------------------------------|-----------------------|
    | WaitTime       | Time (in seconds) to wait for pods   | 60                    |
    | build          | Whether to build the Docker image(Has to be run the first time)    | $true                 |
    | SeedUrl        | Initial URL to crawl                 | "https://example.com" |
    | WorkerCount    | Number of concurrent workers         | 100                   |
    | MaxUrls        | Maximum number of URLs to scrape     | 10000                 |
    | StorageDriver  | Backend storage driver (e.g. postgres) | "postgres"         |

## Custom Drivers
### Adding a New Storage Driver
1. Implement the `StorageDriver` interface defined in` pkg/storage/models/models.go`
2. Create a new driver package in `pkg/database/drivers/customdriver.go`
3. Register the driver in the DriverFactory method in `pkg/storage/storage.go`.

```sh
func DriverFactory(driverType string) (models.StorageDriver, error) {
    switch driverType {
    case "postgres":
        return postgresdriver.NewPostgresDriver()
    // Add your custom driver here 
    default:
        return nil, errors.New("unsupported storage driver type")
    }
}
```
