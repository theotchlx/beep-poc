package api

// Public, unauthenticated API methods.

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// Public API interface, struct, constructor and methods.

type PublicAPI struct {
	server *echo.Echo
}

func InitPublicAPI() *PublicAPI {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	return &PublicAPI{
		server: e,
	}
}

// getWellKnownConfig returns the body of the /auth-well-known-config endpoint.
func (api *PublicAPI) getWellKnownConfig(c echo.Context) error {
	log.Println("getWellKnownConfig endpoint hit")
	resp, err := http.Get("http://localhost:7080/realms/beep-poc/.well-known/openid-configuration")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch well-known configuration"})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.JSON(resp.StatusCode, map[string]string{"error": "unexpected status code from well-known endpoint"})
	}

	var config map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to parse well-known configuration"})
	}

	return c.JSON(http.StatusOK, config)
}
