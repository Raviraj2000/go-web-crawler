// pkg/storage/storage.go
package models

// PageData represents the data structure for crawled pages.
type PageData struct {
	URL         string
	Title       string
	Description string
}

// StorageDriver defines the methods required by all storage drivers.
type StorageDriver interface {
	Save(data PageData) error
	Close() error
}
