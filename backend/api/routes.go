package api

// API routes definition.

import "github.com/labstack/echo/v4/middleware"

func (api *MessageAPI) RegisterRoutes() {
	// Middleware
	api.server.Use(middleware.Logger())
	api.server.Use(middleware.Recover())

	// API routes
	api.server.POST("/messages", api.createMessage) // Create or update a message
	api.server.GET("/messages", api.getPaginatedMessages) // Get messages with pagination
	api.server.GET("/messages/:id", api.getMessage) // Get a message by ID
	api.server.POST("/messages/:id", api.updateMessage) // Update a message by its ID
	api.server.GET("/search/messages", api.searchMessages) // Search messages
}

func (api *MessageAPI) Start(port string) {
	api.RegisterRoutes()
	api.server.Logger.Fatal(api.server.Start(port))
}
