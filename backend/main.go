package main

import (
	"beep-poc-backend/api"
	"beep-poc-backend/repository/elastic"
	"beep-poc-backend/service"
)

func main() {
	// Initialize repositories, services, broker...

	repository := elastic.NewInMemoryMessageRepository() // Init Messages repository
	service := service.InitMessageService(repository) // Init Messages/Gateway service API functions.
	api := api.InitMessageAPI(service) // Init HTTP and broker APIs with the service.

	// Register API routes and start server.
	api.Start(":8080")
}
