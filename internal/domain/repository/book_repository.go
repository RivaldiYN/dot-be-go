package repository

import (
	"errors"

	"dot-be-go/internal/domain/entity"

	"gorm.io/gorm"
)

// BookRepository interface for book operations
type BookRepository interface {
	Create(book *entity.Book) error
	FindAll(userID uint) ([]entity.Book, error)
	FindByID(id uint, userID uint) (*entity.Book, error)
	Update(book *entity.Book) error
	Delete(id uint, userID uint) error
	FindByCategory(categoryID uint) ([]entity.Book, error)
}

// bookRepository implements BookRepository
type bookRepository struct {
	db *gorm.DB
}

// NewBookRepository creates a new book repository
func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db}
}

// Create creates a new book
func (r *bookRepository) Create(book *entity.Book) error {
	return r.db.Create(book).Error
}

// FindAll returns all books for a user
func (r *bookRepository) FindAll(userID uint) ([]entity.Book, error) {
	var books []entity.Book
	err := r.db.Where("user_id = ?", userID).
		Preload("Categories").
		Find(&books).Error
	return books, err
}

// FindByID finds a book by ID for a specific user
func (r *bookRepository) FindByID(id uint, userID uint) (*entity.Book, error) {
	var book entity.Book
	err := r.db.Where("id = ? AND user_id = ?", id, userID).
		Preload("Categories").
		First(&book).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("book not found")
		}
		return nil, err
	}
	return &book, nil
}

// Update updates a book
func (r *bookRepository) Update(book *entity.Book) error {
	// First clear existing category associations
	err := r.db.Model(book).Association("Categories").Clear()
	if err != nil {
		return err
	}

	// Then save the book with new associations
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(book).Error
}

// Delete deletes a book
func (r *bookRepository) Delete(id uint, userID uint) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&entity.Book{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("book not found or you don't have permission")
	}
	return nil
}

// FindByCategory finds books by category ID
func (r *bookRepository) FindByCategory(categoryID uint) ([]entity.Book, error) {
	var books []entity.Book
	err := r.db.Joins("JOIN book_categories ON books.id = book_categories.book_id").
		Where("book_categories.category_id = ?", categoryID).
		Preload("Categories").
		Find(&books).Error
	return books, err
}
