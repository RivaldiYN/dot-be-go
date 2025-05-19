package main

import (
	"fmt"
	"strconv"

	"dot-be-go/config"
	"dot-be-go/internal/app/api/handlers"
	"dot-be-go/internal/app/api/routes"
	"dot-be-go/internal/domain/entity"
	"dot-be-go/internal/domain/repository"
	"dot-be-go/internal/service"
	"dot-be-go/pkg/hash"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg := config.New()

	// Setup database
	db := setupDatabase(cfg)

	// Auto migrate database models
	err := db.AutoMigrate(
		&entity.User{},
		&entity.Category{},
		&entity.Book{},
	)
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	// Create admin user if not exists
	createAdminUser(db)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	bookRepo := repository.NewBookRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.JWTSecretKey, cfg.JWTExpiry)
	categoryService := service.NewCategoryService(categoryRepo)
	bookService := service.NewBookService(bookRepo, categoryRepo)

	// Initialize handlers
	handler := handlers.NewHandler(authService, bookService, categoryService)

	// Setup Echo
	e := echo.New()

	// Setup routes
	routes.SetupRoutes(e, handler, cfg.JWTSecretKey)

	// Start server
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(cfg.AppPort)))
}

func setupDatabase(cfg *config.Config) *gorm.DB {
	var db *gorm.DB
	var err error

	if cfg.DBDriver == "postgres" {
		db, err = gorm.Open(postgres.Open(cfg.DBConnectionString()), &gorm.Config{})
	} else {
		db, err = gorm.Open(mysql.Open(cfg.DBConnectionString()), &gorm.Config{})
	}

	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	return db
}

func createAdminUser(db *gorm.DB) {
	// Check if admin user exists
	var count int64
	db.Model(&entity.User{}).Where("role = ?", "admin").Count(&count)
	if count > 0 {
		return
	}

	// Create admin user if not exists
	hashedPassword, _ := hash.GenerateHash("admin123")
	adminUser := entity.User{
		Name:     "Admin",
		Email:    "admin@example.com",
		Password: hashedPassword,
		Role:     "admin",
	}

	db.Create(&adminUser)
	fmt.Println("Admin user created with email: admin@example.com and password: admin123")
}
