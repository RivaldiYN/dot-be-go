package handlers

import (
	"net/http"
	"strconv"

	"dot-be-go/internal/service"

	"github.com/labstack/echo/v4"
)

// CreateBook creates a new book
func (h *Handler) CreateBook(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	req := new(service.BookRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	book, err := h.BookService.Create(userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, book)
}

// GetAllBooks returns all books for a user
func (h *Handler) GetAllBooks(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	books, err := h.BookService.GetAll(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, books)
}

// GetBookByID returns a book by ID for a user
func (h *Handler) GetBookByID(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid book ID")
	}

	book, err := h.BookService.GetByID(uint(id), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, book)
}

// UpdateBook updates a book
func (h *Handler) UpdateBook(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid book ID")
	}

	req := new(service.BookRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	book, err := h.BookService.Update(uint(id), userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, book)
}

// DeleteBook deletes a book
func (h *Handler) DeleteBook(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid book ID")
	}

	if err := h.BookService.Delete(uint(id), userID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

// GetBooksByCategory returns books by category
func (h *Handler) GetBooksByCategory(c echo.Context) error {
	categoryIDParam := c.Param("categoryId")
	categoryID, err := strconv.ParseUint(categoryIDParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid category ID")
	}

	books, err := h.BookService.GetByCategory(uint(categoryID))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, books)
}
