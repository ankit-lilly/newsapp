package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ankit-lilly/newsmaxxing/pkg/auth"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
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
		NewClaimsFunc: func(c *echo.Context) jwt.Claims {
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

func (m *AuthMiddleware) jwtSuccessHandler(c *echo.Context) error {
	token, err := echo.ContextGet[*jwt.Token](c, "user")
	if err != nil {
		c.Logger().Error("token not found in context", "error", err)
		return err
	}

	claims, ok := token.Claims.(*auth.JwtClaims)
	if !ok {
		err := fmt.Errorf("unexpected JWT claims type %T", token.Claims)
		c.Logger().Error("failed to get claims from token", "error", err)
		return err
	}

	c.Set("userId", claims.Id)
	c.Set("userName", claims.Username)
	c.Set("isAuthorized", true)
	c.Set("currentPath", c.Request().URL.Path)

	c.Logger().Info("user is authorized", "path", c.Path(), "user_id", claims.Id)

	switch c.Path() {
	case "/login", "/register":
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	return nil
}

func (m *AuthMiddleware) jwtErrorHandler(c *echo.Context, err error) error {
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
