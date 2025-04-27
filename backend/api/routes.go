package api

// API routes definition.

import "github.com/labstack/echo/v4/middleware"

func (api *MessageAPI) RegisterRoutes() {
	// Middleware
	api.server.Use(middleware.Logger())
	api.server.Use(middleware.Recover())

	// API routes
	endpoints := api.server.Group("/api")
	endpoints.GET("/messages", api.getAllMessages)
	endpoints.GET("/messages/:id", api.getMessage)
	endpoints.POST("/messages", api.createMessage)
}

func (api *MessageAPI) Start(port string) {
	api.RegisterRoutes()
	api.server.Logger.Fatal(api.server.Start(port))
}
