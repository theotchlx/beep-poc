package authn

import (
	"context"
	"log"
	"net/http"

	"github.com/coreos/go-oidc"
	"github.com/labstack/echo/v4"
)

// Config holds Keycloak settings.
type Config struct {
	IssuerURL string // dev would be "http://localhost:7080/realms/beep-poc"
	ClientID  string // dev would be "beep-poc-front"
}

// AuthMiddleware wraps an OIDC token verifier.
type AuthMiddleware struct {
	verifier *oidc.IDTokenVerifier
}

// NewAuthMiddleware creates a new AuthMiddleware from Keycloak config.
func NewAuthMiddleware(cfg Config) (*AuthMiddleware, error) {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, cfg.IssuerURL)
	if err != nil {
		return nil, err
	}

	oidcConfig := &oidc.Config{
		ClientID:          cfg.ClientID,
		SkipClientIDCheck: true, // Skip strict audience validation
	}
	verifier := provider.Verifier(oidcConfig)

	return &AuthMiddleware{verifier: verifier}, nil
}

// MiddlewareFunc returns an Echo middleware that enforces a valid Bearer JWT.
func (mw *AuthMiddleware) MiddlewareFunc() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			log.Println("Middleware executed")
			token := c.Request().Header.Get("Authorization")
			if token == "" {
				log.Println("No Authorization header provided")
                return echo.NewHTTPError(http.StatusUnauthorized, "Missing token")
			}

			// Remove "Bearer " prefix if present
			if len(token) > 7 && token[:7] == "Bearer " {
				token = token[7:]
			}

			// Validate the token
			claims, err := mw.ValidateToken(token, c)
			if err != nil {
				log.Printf("Token validation failed: %v", err)
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}

			log.Printf("Token validated successfully: %v", claims)
			return next(c)
		}
	}
}

// ValidateToken verifies the token, extracts claims, and stores them in context.
func (a *AuthMiddleware) ValidateToken(token string, c echo.Context) (string, error) {
	ctx := c.Request().Context()
	idToken, err := a.verifier.Verify(ctx, token)
	if err != nil {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
	}

	var claims struct {
		Subject     string `json:"sub"`
		Email       string `json:"email"`
		RealmAccess struct {
			Roles []string `json:"roles"`
		} `json:"realm_access"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "failed to parse claims")
	}

	// Expose user info to handlers
	c.Set("userID", claims.Subject)
	c.Set("email", claims.Email)
	c.Set("roles", claims.RealmAccess.Roles)

	return token, nil
}
