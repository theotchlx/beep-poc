package main

import (
	"beep-poc-backend/api"
	"beep-poc-backend/repository/elastic"
	"beep-poc-backend/service"
	"fmt"

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

	message, err := repository.Get("1a") // Test the repository with a sample ID.
	if err != nil {
		log.Fatalf("Error getting message: %s", err)
	}
	fmt.Printf("Message: %v\n", message)

	service := service.InitMessageService(repository) // Init Messages/Gateway service API functions.
	api := api.InitMessageAPI(service) // Init HTTP and broker APIs with the service.

	// Register API routes and start server.
	api.Start(":8080")
}
