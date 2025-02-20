package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ankit-lilly/newsapp/internal/models"
	"github.com/ankit-lilly/newsapp/internal/services"
	"github.com/ankit-lilly/newsapp/internal/templates/components/ui"
	"github.com/ankit-lilly/newsapp/internal/templates/components/users"
	"github.com/ankit-lilly/newsapp/internal/templates/pages"
	"github.com/ankit-lilly/newsapp/pkg/auth"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	*BaseHandler
	userService *services.UserService
	jwtService  *auth.JwtService
}

func NewUserHandler(userService *services.UserService, jwtService *auth.JwtService) *UserHandler {
	return &UserHandler{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (h *UserHandler) LoginView(c echo.Context) error {

	return h.Render(c, RenderProps{
		Title:            "Login",
		Component:        users.Login(),
		WrapperComponent: pages.Index,
	})
}

func (h *UserHandler) RegisterView(c echo.Context) error {
	return h.Render(c, RenderProps{
		Title:            "Register",
		Component:        users.Register(),
		WrapperComponent: pages.Index,
	})
}

func (h *UserHandler) Register(c echo.Context) error {
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return h.View(c, ui.ErrorBlock(err.Error()))
	}

	exists, err := h.userService.UserExists(c.Request().Context(), user.Email)

	if err != nil {
		c.Logger().Error("Error checking if user exists", err)
		return h.View(c, ui.ErrorBlock("Internal server error."))
	}

	if exists {
		return h.View(c, ui.ErrorBlock("User already exists"))
	}

	id, err := h.userService.Create(c.Request().Context(), &user)

	if err != nil {
		return h.View(c, ui.ErrorBlock(err.Error()))
	}

	err = h.jwtService.GenerateTokenAndSetCookie(auth.User{ID: id, Username: user.Username}, c)

	if err != nil {
		c.Logger().Error("Error generating token", err)
		return h.View(c, ui.ErrorBlock("Internal server error."))
	}

	c.Response().Header().Set("HX-Redirect", fmt.Sprintf("/?v=%d", time.Now().Unix()))
	return nil
}

func (h *UserHandler) Login(c echo.Context) error {
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	u, err := h.userService.ValidateUserCredentials(c.Request().Context(), user.Email, user.Password)

	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			return h.View(c, ui.ErrorBlock("Invalid credentials"))
		}
		return h.View(c, ui.ErrorBlock(err.Error()))
	}

	err = h.jwtService.GenerateTokenAndSetCookie(auth.User{ID: u.ID, Username: u.Username}, c)

	if err != nil {
		c.Logger().Error("Error generating token", err)
		return h.View(c, ui.ErrorBlock("Internal server error."))
	}

	c.Response().Header().Set("HX-Redirect", fmt.Sprintf("/?v=%d", time.Now().Unix()))
	return nil
}

func (h *UserHandler) Logout(c echo.Context) error {
	h.jwtService.Logout(c)

	c.Response().Header().Set("HX-Redirect", fmt.Sprintf("/login?v=%d", time.Now().Unix()))
	return nil

}
