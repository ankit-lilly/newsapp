package handlers

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type AuthService interface {
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

	return cmp.Render(c.Request().Context(), c.Response().Writer)
}

func (a *AuthHandler) LoginHandler(c echo.Context) error {
	return nil
}

func (a *AuthHandler) RegisterUser(c echo.Context) error {
	return nil
}

func (a *AuthHandler) LoginUser(c echo.Context) error {
	return nil
}
