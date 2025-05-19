package main

import (
	"beep-poc-backend/api"
	"beep-poc-backend/repository/elastic"
	"beep-poc-backend/service"

	"log"
	"os"

	"github.com/elastic/go-elasticsearch/v9"
)

func main() {
	// Initialize repositories, services, api...
	client, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{os.Getenv("ES_ADDRESS")},
		Username:  os.Getenv("ELASTICSEARCH_USERNAME"),
		Password:  os.Getenv("ELASTICSEARCH_PASSWORD"),
	})
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	repository := elastic.NewMessageRepository(client) // Init Elasticsearch Messages repository
	service := service.InitMessageService(repository)  // Init Messages/Gateway service API functions.
	messApi := api.InitMessageAPI(service)             // Init HTTP APIs with the service.
	pubApi := api.InitPublicAPI()                      // Init HTTP APIs with the service.

	// Register API routes and start server.
	api.Start(messApi, pubApi, ":8080")
}
