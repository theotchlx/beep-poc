package main

import (
	"beep-poc-backend/api"
	"beep-poc-backend/repository/elastic"
	"beep-poc-backend/service"

	"log"
	"github.com/elastic/go-elasticsearch/v9"
)

func main() {
	// Initialize repositories, services, api...
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	repository := elastic.NewMessageRepository(client) // Init Elasticsearch Messages repository

	service := service.InitMessageService(repository) // Init Messages/Gateway service API functions.
	api := api.InitMessageAPI(service) // Init HTTP and broker APIs with the service.

	// Register API routes and start server.
	api.Start(":8080")
}
