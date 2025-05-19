package api

import (
	"github.com/labstack/echo/v4/middleware"
)

// API routes definition.

func (api *MessageAPI) RegisterMessageRoutes() {
	// Echo middlewares
	api.server.Use(middleware.Logger())
	api.server.Use(middleware.Recover())

	// API routes
	api.server.POST("/messages", api.createMessage)        // Create or update a message
	api.server.GET("/messages", api.getPaginatedMessages)  // Get messages with pagination
	api.server.GET("/messages/:id", api.getMessage)        // Get a message by ID
	api.server.POST("/messages/:id", api.updateMessage)    // Update a message by its ID
	api.server.GET("/search/messages", api.searchMessages) // Search messages
}

func (api *PublicAPI) RegisterPublicRoutes() {
	// Echo middlewares
	api.server.Use(middleware.Logger())
	api.server.Use(middleware.Recover())

	// Routes to manage authentication
	api.server.GET("/public/auth-well-known-config", api.getWellKnownConfig) // Get realm OIDC config
}

func Start(messApi *MessageAPI, pubApi *PublicAPI, port string) {
	messApi.RegisterMessageRoutes()
	messApi.server.Logger.Fatal(messApi.server.Start(port))

	pubApi.RegisterPublicRoutes()
	pubApi.server.Logger.Fatal(pubApi.server.Start(port))
}
