package api

import (
	authn "beep-poc-backend/middlewares/authentication"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// API routes definition.

func (api *MessageAPI) RegisterMessageRoutes(group *echo.Group) {
	// Protected API routes
	group.POST("/messages", api.createMessage)        // Create or update a message
	group.GET("/messages", api.getPaginatedMessages)  // Get messages with pagination
	group.GET("/messages/:id", api.getMessage)        // Get a message by ID
	group.POST("/messages/:id", api.updateMessage)    // Update a message by its ID
	group.GET("/search/messages", api.searchMessages) // Search messages
}

func (api *PublicAPI) RegisterPublicRoutes(group *echo.Group) {
	// Routes to manage authentication
	group.GET("/auth-well-known-config", api.getWellKnownConfig) // Get realm OIDC config
}

func Start(messApi *MessageAPI, pubApi *PublicAPI, port string) {
	e := echo.New()

	// Echo middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Initialize Keycloak auth middleware
	keycloakCfg := authn.Config{
		IssuerURL: "http://localhost:7080/realms/beep-poc",
		ClientID:  "beep-poc-front",
	}
	authMw, err := authn.NewAuthMiddleware(keycloakCfg)
	if err != nil {
		log.Fatalf("failed to init Keycloak auth: %v", err)
	}

	// Public routes (no authentication)
	publicGroup := e.Group("/public")
	pubApi.RegisterPublicRoutes(publicGroup)

	// Protected routes (with authentication)
	protectedGroup := e.Group("")
	protectedGroup.Use(authMw.MiddlewareFunc())
	messApi.RegisterMessageRoutes(protectedGroup)

	// Start the server
	e.Logger.Fatal(e.Start(port))
}
