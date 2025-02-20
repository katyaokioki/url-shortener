//

package main

import (
	"flag"
	"log"
	"net/http"

	"url-shortener/internal/handler"
	"url-shortener/internal/storage"
)

func main() {
	// Флаг для выбора хранилища
	storageType := flag.String("storage_type", "memory", "Storage type: memory or postgres")
	postgresConn := flag.String("postgres_dsn", "host=localhost user=postgres password=pass dbname=url_shortener sslmode=disable", "PostgreSQL DSN")

	flag.Parse()

	var store storage.Storage
	var err error

	if *storageType == "postgres" {
		store, err = storage.NewPostgresStorage(*postgresConn)
		if err != nil {
			log.Fatalf("Failed to connect to Postgres: %v", err)
		}
		log.Println("Using PostgreSQL as storage")
	} else {
		store = storage.NewInMemoryStorage()
		log.Println("Using in-memory storage")
	}

	server := handler.NewServer(store)

	http.HandleFunc("/shorten", server.HandleShorten)
	http.HandleFunc("/expand", server.HandleExpand)

	log.Println("HTTP Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
