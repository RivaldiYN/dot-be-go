package handlers

import (
	"dot-be-go/internal/service"
)

// Handler contains all handlers for API endpoints
type Handler struct {
	AuthService     service.AuthService
	BookService     service.BookService
	CategoryService service.CategoryService
}

// NewHandler creates a new handler instance
func NewHandler(
	authService service.AuthService,
	bookService service.BookService,
	categoryService service.CategoryService,
) *Handler {
	return &Handler{
		AuthService:     authService,
		BookService:     bookService,
		CategoryService: categoryService,
	}
}
