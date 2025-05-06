package authn

import (
    "context"
    "net/http"

    "github.com/coreos/go-oidc"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
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

    oidcConfig := &oidc.Config{ClientID: cfg.ClientID}
    verifier := provider.Verifier(oidcConfig)
    return &AuthMiddleware{verifier: verifier}, nil
}

// MiddlewareFunc returns an Echo middleware that enforces a valid Bearer JWT.
func (a *AuthMiddleware) MiddlewareFunc() echo.MiddlewareFunc {
    return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
        KeyLookup:  "header:Authorization",
        AuthScheme: "Bearer",
        Validator: func(token string, c echo.Context) (bool, error) {
            _, err := a.validateToken(token, c)
            return err == nil, err
        },
    })
}

// validateToken verifies the token, extracts claims, and stores them in context.
func (a *AuthMiddleware) validateToken(token string, c echo.Context) (string, error) {
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

