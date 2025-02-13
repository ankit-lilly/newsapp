package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/ankit-lilly/newsapp/internal/auth/respository"
	"github.com/ankit-lilly/newsapp/internal/auth/views"
	"github.com/ankit-lilly/newsapp/pkg/auth"
	"github.com/ankit-lilly/newsapp/pkg/shared"
	"github.com/ankit-lilly/newsapp/pkg/views/components"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type AuthService interface {
	CreateUser(user repository.User) (*repository.User, error)
	LoginUser(email, pasword string) (*repository.User, error)
	GetUserByEmail(email string) (*repository.User, error)
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
	c.Response().Header().Set("Cache-Control", "private, max-age=86400, stale-while-revalidate=30")
	if isAuthorized {
		c.Response().Header().Set("Hx-Redirect", "/")
		return c.Redirect(http.StatusOK, "/")
	}
	si := views.ShowLogin("| Login", isAuthorized, shared.Categories, views.Login())
	return a.View(c, si)
}

func (a *AuthHandler) RegisterHandler(c echo.Context) error {

	isAuthorized := c.Get("isAuthorized").(bool)

	c.Response().Header().Set("Cache-Control", "private, max-age=86400, stale-while-revalidate=30")
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
		c.Response().Header().Set("HX-Retarget", "#validation-error-block")
		return a.View(c, components.ErrorBlock("Email/Password is required."))
	}

	user, err := a.AuthService.LoginUser(email, password)

	if err != nil {
		msg := err.Error()

		if strings.Contains(msg, "not found") || strings.Contains(msg, "Invalid password") {
			c.Response().Header().Set("HX-Retarget", "#validation-error-block")
			return a.View(c, components.ErrorBlock("Invalid email/password combination"))
		}

		c.Response().Header().Set("HX-Retarget", "#validation-error-block")
		return a.View(c, components.ErrorBlock("Something went wrong"))
	}

	duration := 1 * time.Hour
	tokenErr := auth.GenerateTokensAndSetCookies(user, duration, c)

	if tokenErr != nil {
		log.Error("Couldn't set cookies", tokenErr)
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, "Internal server error")
	}

	c.Response().Header().Set("HX-Redirect", "/?v="+uuid.New().String())
	return c.Redirect(http.StatusOK, "/")
}

func (a *AuthHandler) RegisterUser(c echo.Context) error {

	name, email, password := c.FormValue("name"), c.FormValue("email"), c.FormValue("password")

	if name == "" || email == "" || password == "" {

		c.Response().Header().Set("HX-Retarget", "#validation-error-block")
		return a.View(c, components.ErrorBlock("Missing required fields"))
	}

	user := repository.User{Name: name, Email: email, Password: password}

	userExists, err := a.AuthService.GetUserByEmail(email)

	if userExists != nil {
		c.Response().Header().Set("HX-Retarget", "#validation-error-block")
		return a.View(c, components.ErrorBlock("User with this email already exists"))
	}

	createdUser, err := a.AuthService.CreateUser(user)

	if err != nil {
		log.Error("Couldn't add user", err)
		c.Response().Header().Set("HX-Retarget", "#validation-error-block")
		return a.View(c, components.ErrorBlock("Something went wrong. Please try again"))
	}

	duration := 1 * time.Hour
	tokenErr := auth.GenerateTokensAndSetCookies(createdUser, duration, c)

	if tokenErr != nil {
		log.Error("Couldn't set cookies", tokenErr)
		c.Response().Header().Set("HX-Retarget", "#validation-error-block")
		return a.View(c, components.ErrorBlock("Something went wrong."))
	}

	c.Response().Header().Set("HX-Redirect", "/?v="+uuid.New().String())
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

	c.Response().Header().Set("HX-Redirect", "/?v="+uuid.New().String())
	return c.Redirect(http.StatusOK, "/auth/login")
}
