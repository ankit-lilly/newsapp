package auth

import (
	"database/sql"

	handlers "github.com/ankibahuguna/newsapp/internal/auth/handlers"
	repository "github.com/ankibahuguna/newsapp/internal/auth/respository"
	services "github.com/ankibahuguna/newsapp/internal/auth/services"
	"github.com/labstack/echo/v4"
)

func Routes(a *echo.Group, DB *sql.DB) {
	authRepository := repository.NewRepository(DB)
	authService := services.NewAuthService(authRepository)
	authHandler := handlers.New(authService)

	a.GET("/login", authHandler.LoginHandler)
	a.POST("/login", authHandler.LoginUser)

	a.GET("/register", authHandler.RegisterHandler)
	a.POST("/register", authHandler.RegisterUser)

	a.DELETE("/logout", authHandler.LogoutUser)
}
