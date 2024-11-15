package storage

// PageData represents the structure of the crawled data.
type PageData struct {
	URL        string
	Title      string
	Desciption string
}

// StorageDriver defines the interface for different storage backends.
type StorageDriver interface {
	Configure() error // Allow the driver to self-configure (e.g., load environment variables, configs)
	Save(data PageData) error
	Close() error
}
