package routes

import (
	"database/sql"

	"github.com/ankit-lilly/newsapp/internal/handlers"
	"github.com/ankit-lilly/newsapp/internal/repositories"
	"github.com/ankit-lilly/newsapp/internal/services"
	"github.com/ankit-lilly/newsapp/internal/services/llm"
	"github.com/ankit-lilly/newsapp/internal/services/providers"
	"github.com/ankit-lilly/newsapp/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/olahol/melody"
)

// RegisterRoutes wires up the application's routes and dependencies.
func RegisterRoutes(e *echo.Echo, db *sql.DB, llmHandler *llm.LLMHandler) {

	articleRepository := repositories.NewArticleRepository(db)
	articleService := services.NewArticleService(articleRepository, llmHandler, providers.Registry["thehindu"])
	articleHandler := handlers.NewArticleHandler(articleService)

	e.GET("/", articleHandler.List).Name = "homePage"
	e.GET("/news/:portal", articleHandler.List)
	e.GET("/news/:portal/:category", articleHandler.ListByCategory)
	e.GET("/articles/:portal/:id", articleHandler.GetArticleByID)
	e.GET("/articles/:portal/:id/summarise", articleHandler.GetArticleSummary)
	e.POST("/articles/:portal/:id/favorites", articleHandler.CreateFavoriteArticle)
	e.DELETE("/articles/:id", articleHandler.DeleteFavoriteArticle)
	e.GET("/favorites", articleHandler.ListFavoriteArticles)

	// ----- User Routes -----
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService, auth.NewJwtService())

	e.GET("/login", userHandler.LoginView).Name = "loginPage"
	e.POST("/login", userHandler.Login)
	e.GET("/register", userHandler.RegisterView).Name = "registerPage"
	e.POST("/register", userHandler.Register)
	e.DELETE("/logout", userHandler.Logout)

	// ----- Chat Routes (WebSocket) -----
	ws := melody.New()
	chatHandler := handlers.NewChatHandler(articleService, ws)
	e.GET("/articles/:portal/:id/chat", chatHandler.Chat)
	ws.HandleMessage(chatHandler.HandleChatMessage)
	ws.HandleConnect(chatHandler.HandleConnect)
	ws.HandleDisconnect(chatHandler.HandleDisconnect)
}
