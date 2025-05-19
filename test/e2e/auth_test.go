package e2e

import (
	"bytes"
	"dot-be-go/config"
	"dot-be-go/internal/app/api/handlers"
	"dot-be-go/internal/app/api/routes"
	"dot-be-go/internal/domain/entity"
	"dot-be-go/internal/domain/repository"
	"dot-be-go/internal/service"
	"dot-be-go/pkg/hash"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func setupTestEnvironment(t *testing.T) (*echo.Echo, *gorm.DB, *handlers.Handler) {
	cfg := config.New()
	cfg.DBName = "bookdb"
	db, err := gorm.Open(postgres.Open(cfg.DBConnectionString()), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	db.Exec("TRUNCATE TABLE users CASCADE")
	db.Exec("TRUNCATE TABLE categories CASCADE")
	db.Exec("TRUNCATE TABLE books CASCADE")

	err = db.AutoMigrate(
		&entity.User{},
		&entity.Category{},
		&entity.Book{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	hashedPassword, _ := hash.GenerateHash("admin123")
	adminUser := entity.User{
		Name:     "Admin",
		Email:    "admin@example.com",
		Password: hashedPassword,
		Role:     "admin",
	}
	result := db.Create(&adminUser)
	if result.Error != nil {
		t.Fatalf("Failed to create test admin user: %v", result.Error)
	}

	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	bookRepo := repository.NewBookRepository(db)

	authService := service.NewAuthService(userRepo, cfg.JWTSecretKey, cfg.JWTExpiry)
	categoryService := service.NewCategoryService(categoryRepo)
	bookService := service.NewBookService(bookRepo, categoryRepo)

	handler := handlers.NewHandler(authService, bookService, categoryService)

	e := echo.New()

	routes.SetupRoutes(e, handler, cfg.JWTSecretKey)

	return e, db, handler
}

func TestLoginSuccess(t *testing.T) {
	e, db, _ := setupTestEnvironment(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	loginReq := loginRequest{
		Email:    "admin@example.com",
		Password: "admin123",
	}
	jsonBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response loginResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Token)
}

func TestLoginFailure_InvalidCredentials(t *testing.T) {
	e, db, _ := setupTestEnvironment(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	loginReq := loginRequest{
		Email:    "admin@example.com",
		Password: "wrongpassword",
	}
	jsonBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestLoginFailure_InvalidEmail(t *testing.T) {
	e, db, _ := setupTestEnvironment(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()
	loginReq := loginRequest{
		Email:    "nonexistent@example.com",
		Password: "admin123",
	}
	jsonBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestLoginFailure_InvalidJSON(t *testing.T) {
	e, db, _ := setupTestEnvironment(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()
	invalidJSON := []byte(`{"email": "admin@example.com", "password":}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(invalidJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
