package service

import (
	"errors"

	"dot-be-go/internal/domain/entity"
	"dot-be-go/internal/domain/repository"
)

// BookRequest represents book request data
type BookRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Author      string `json:"author" validate:"required,min=1,max=100"`
	ISBN        string `json:"isbn" validate:"required,min=10,max=20"`
	PublishYear int    `json:"publish_year" validate:"required,min=1000,max=9999"`
	Description string `json:"description" validate:"max=1000"`
	CategoryIDs []uint `json:"category_ids" validate:"dive,min=1"`
}

// BookService handles book operations
type BookService interface {
	Create(userID uint, req *BookRequest) (*entity.Book, error)
	GetAll(userID uint) ([]entity.Book, error)
	GetByID(id uint, userID uint) (*entity.Book, error)
	Update(id uint, userID uint, req *BookRequest) (*entity.Book, error)
	Delete(id uint, userID uint) error
	GetByCategory(categoryID uint) ([]entity.Book, error)
}

type bookService struct {
	bookRepo     repository.BookRepository
	categoryRepo repository.CategoryRepository
}

// NewBookService creates a new book service
func NewBookService(bookRepo repository.BookRepository, categoryRepo repository.CategoryRepository) BookService {
	return &bookService{
		bookRepo:     bookRepo,
		categoryRepo: categoryRepo,
	}
}

// Create creates a new book
func (s *bookService) Create(userID uint, req *BookRequest) (*entity.Book, error) {
	book := &entity.Book{
		Title:       req.Title,
		Author:      req.Author,
		ISBN:        req.ISBN,
		PublishYear: req.PublishYear,
		Description: req.Description,
		UserID:      userID,
		Categories:  []entity.Category{},
	}

	// Add categories
	if len(req.CategoryIDs) > 0 {
		for _, categoryID := range req.CategoryIDs {
			category, err := s.categoryRepo.FindByID(categoryID)
			if err != nil {
				return nil, errors.New("category not found: " + err.Error())
			}
			book.Categories = append(book.Categories, *category)
		}
	}

	if err := s.bookRepo.Create(book); err != nil {
		return nil, err
	}

	return book, nil
}

// GetAll returns all books for a user
func (s *bookService) GetAll(userID uint) ([]entity.Book, error) {
	return s.bookRepo.FindAll(userID)
}

// GetByID returns a book by ID for a specific user
func (s *bookService) GetByID(id uint, userID uint) (*entity.Book, error) {
	return s.bookRepo.FindByID(id, userID)
}

// Update updates a book
func (s *bookService) Update(id uint, userID uint, req *BookRequest) (*entity.Book, error) {
	book, err := s.bookRepo.FindByID(id, userID)
	if err != nil {
		return nil, err
	}

	book.Title = req.Title
	book.Author = req.Author
	book.ISBN = req.ISBN
	book.PublishYear = req.PublishYear
	book.Description = req.Description

	// Replace categories
	book.Categories = []entity.Category{}
	if len(req.CategoryIDs) > 0 {
		for _, categoryID := range req.CategoryIDs {
			category, err := s.categoryRepo.FindByID(categoryID)
			if err != nil {
				return nil, errors.New("category not found: " + err.Error())
			}
			book.Categories = append(book.Categories, *category)
		}
	}

	if err := s.bookRepo.Update(book); err != nil {
		return nil, err
	}

	return book, nil
}

// Delete deletes a book
func (s *bookService) Delete(id uint, userID uint) error {
	return s.bookRepo.Delete(id, userID)
}

// GetByCategory returns books by category ID
func (s *bookService) GetByCategory(categoryID uint) ([]entity.Book, error) {
	return s.bookRepo.FindByCategory(categoryID)
}
