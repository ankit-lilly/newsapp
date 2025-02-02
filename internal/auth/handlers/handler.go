package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/ankibahuguna/newsapp/internal/auth/respository"
	"github.com/ankibahuguna/newsapp/internal/auth/views"
	"github.com/ankibahuguna/newsapp/pkg/auth"
	"github.com/ankibahuguna/newsapp/pkg/errors"
	"github.com/ankibahuguna/newsapp/pkg/shared"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type AuthService interface {
	CreateUser(user repository.User) (*repository.User, error)
	LoginUser(email, pasword string) (*repository.User, error)
}

func New(a AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: a,
	}
}

type AuthHandler struct {
	AuthService AuthService
}

func (a *AuthHandler) View(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	ctx := context.WithValue(context.Background(), "currentPath", c.Request().URL.Path)
	return cmp.Render(ctx, c.Response().Writer)
}

func (a *AuthHandler) LoginHandler(c echo.Context) error {
	isAuthorized := c.Get("isAuthorized").(bool)
	if isAuthorized {
		c.Response().Header().Set("Hx-Redirect", "/")
		return c.Redirect(http.StatusOK, "/")
	}
	si := views.ShowLogin("| Login", isAuthorized, shared.Categories, views.Login())
	return a.View(c, si)
}

func (a *AuthHandler) RegisterHandler(c echo.Context) error {

	isAuthorized := c.Get("isAuthorized").(bool)

	if isAuthorized {
		c.Response().Header().Set("Hx-Redirect", "/")
		return c.Redirect(http.StatusOK, "")
	}

	si := views.ShowRegister("| Register", isAuthorized, shared.Categories, views.Register())
	return a.View(c, si)
}

func (a *AuthHandler) LoginUser(c echo.Context) error {

	email, password := c.FormValue("email"), c.FormValue("password")

	if email == "" || password == "" {
		return errors.NewApiError("email and passwords are required is required", http.StatusBadRequest)
	}

	user, err := a.AuthService.LoginUser(email, password)

	if err != nil {
		msg := err.Error()

		if strings.Contains(msg, "not found") || strings.Contains(msg, "Invalid password") {
			return echo.NewHTTPError(echo.ErrUnauthorized.Code, "Invalid email/password combination")
		}
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, "Something went wrong")
	}

	duration := 1 * time.Hour
	tokenErr := auth.GenerateTokensAndSetCookies(user, duration, c)

	if tokenErr != nil {
		log.Error("Couldn't set cookies", tokenErr)
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, "Internal server error")
	}

	c.Response().Header().Set("Hx-Redirect", "/")
	return c.Redirect(http.StatusOK, "/")
}

func (a *AuthHandler) RegisterUser(c echo.Context) error {

	name, email, password := c.FormValue("name"), c.FormValue("email"), c.FormValue("password")

	if name == "" || email == "" || password == "" {
		return errors.NewApiError("Item name is required", http.StatusBadRequest)
	}

	user := repository.User{Name: name, Email: email, Password: password}

	createdUser, err := a.AuthService.CreateUser(user)

	if err != nil {
		log.Error("Couldn't add user", err)
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, "Internal server error")
	}

	duration := 1 * time.Hour
	tokenErr := auth.GenerateTokensAndSetCookies(createdUser, duration, c)

	if tokenErr != nil {
		log.Error("Couldn't set cookies", tokenErr)
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, "Internal server error")
	}

	c.Response().Header().Set("Hx-Redirect", "/")
	return c.Redirect(http.StatusOK, "/")
}

func (a *AuthHandler) LogoutUser(c echo.Context) error {

	isAuthorized := c.Get("isAuthorized").(bool)

	if !isAuthorized {
		c.Response().Header().Set("Hx-Redirect", "/auth/login")
		return c.Redirect(http.StatusOK, "/auth/login")
	}

	duration := -1 * time.Hour
	tokenErr := auth.GenerateTokensAndSetCookies(&repository.User{}, duration, c)

	if tokenErr != nil {
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, "Internal server error")
	}

	c.Response().Header().Set("Hx-Redirect", "/")
	return c.Redirect(http.StatusOK, "/auth/login")
}
