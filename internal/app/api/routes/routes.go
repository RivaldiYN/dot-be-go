package routes

import (
	"dot-be-go/internal/app/api/handlers"
	customMiddleware "dot-be-go/internal/app/api/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SetupRoutes sets up API routes
func SetupRoutes(e *echo.Echo, handler *handlers.Handler, jwtSecret string) {
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	// Public routes
	e.POST("/api/auth/register", handler.Register)
	e.POST("/api/auth/login", handler.Login)

	// Public category routes
	e.GET("/api/categories", handler.GetAllCategories)
	e.GET("/api/categories/:id", handler.GetCategoryByID)
	e.GET("/api/categories/:categoryId/books", handler.GetBooksByCategory)

	// Protected routes
	protected := e.Group("/api")
	protected.Use(customMiddleware.JWTMiddleware(jwtSecret))

	// User routes
	protected.GET("/profile", handler.GetProfile)

	// Book routes
	protected.POST("/books", handler.CreateBook)
	protected.GET("/books", handler.GetAllBooks)
	protected.GET("/books/:id", handler.GetBookByID)
	protected.PUT("/books/:id", handler.UpdateBook)
	protected.DELETE("/books/:id", handler.DeleteBook)

	// Admin routes
	admin := protected.Group("/admin")
	admin.Use(customMiddleware.AdminMiddleware())

	// Admin category management
	admin.POST("/categories", handler.CreateCategory)
	admin.PUT("/categories/:id", handler.UpdateCategory)
	admin.DELETE("/categories/:id", handler.DeleteCategory)
}
