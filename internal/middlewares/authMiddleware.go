package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ankit-lilly/newsapp/pkg/auth"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	jwtService *auth.JwtService
}

func NewAuthMiddleware(jwtService *auth.JwtService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

// JWT returns the JWT middleware configuration
func (m *AuthMiddleware) JWT() echo.MiddlewareFunc {
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.JwtClaims)
		},
		TokenLookup:            fmt.Sprintf("cookie:%s", "access-token"),
		SigningKey:             []byte(m.jwtService.JWTSecret),
		SuccessHandler:         m.jwtSuccessHandler,
		ErrorHandler:           m.jwtErrorHandler,
		ContinueOnIgnoredError: true,
	}

	return echojwt.WithConfig(config)
}

func (m *AuthMiddleware) jwtSuccessHandler(c echo.Context) {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		c.Echo().Logger.Error("Token not found in context")
		return
	}

	claims, ok := token.Claims.(*auth.JwtClaims)
	if !ok {
		c.Echo().Logger.Error("Failed to get claims from token")
		return
	}

	c.Set("userId", claims.Id)
	c.Set("userName", claims.Username)
	c.Set("isAuthorized", true)
	c.Set("currentPath", c.Request().URL.Path)

	c.Logger().Info("User is authorized", c.Path())

	switch c.Path() {
	case "/login", "/register":
		c.Redirect(http.StatusTemporaryRedirect, c.Echo().Reverse("homePage"))
	}
}

func (m *AuthMiddleware) jwtErrorHandler(c echo.Context, err error) error {
	slog.Error("JWT validation failed",
		"error", err,
		"ip", c.RealIP(),
		"path", c.Path(),
		"userAgent", c.Request().UserAgent(),
	)

	c.Set("userId", nil)
	c.Set("userName", nil)
	c.Set("isAuthorized", false)

	return nil
}
