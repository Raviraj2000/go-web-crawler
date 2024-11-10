package storage

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type PageData struct {
	URL        string
	Title      string
	Desciption string
}

var (
	file  *os.File
	mutex sync.Mutex
)

func init() {
	var err error
	file, err = os.OpenFile("scraped-data/results.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
}

func Save(data PageData) {
	mutex.Lock()
	defer mutex.Unlock()

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling data: %v", err)
		return
	}

	_, err = file.WriteString(string(jsonData) + "\n")
	if err != nil {
		log.Printf("Error writing to file: %v", err)
	}
}
