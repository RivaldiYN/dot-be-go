package repository

import (
	"errors"

	"dot-be-go/internal/domain/entity"

	"gorm.io/gorm"
)

// CategoryRepository interface for category operations
type CategoryRepository interface {
	Create(category *entity.Category) error
	FindAll() ([]entity.Category, error)
	FindByID(id uint) (*entity.Category, error)
	Update(category *entity.Category) error
	Delete(id uint) error
}

// categoryRepository implements CategoryRepository
type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db}
}

// Create creates a new category
func (r *categoryRepository) Create(category *entity.Category) error {
	return r.db.Create(category).Error
}

// FindAll returns all categories
func (r *categoryRepository) FindAll() ([]entity.Category, error) {
	var categories []entity.Category
	err := r.db.Find(&categories).Error
	return categories, err
}

// FindByID finds a category by ID
func (r *categoryRepository) FindByID(id uint) (*entity.Category, error) {
	var category entity.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return &category, nil
}

// Update updates a category
func (r *categoryRepository) Update(category *entity.Category) error {
	return r.db.Save(category).Error
}

// Delete deletes a category
func (r *categoryRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Category{}, id).Error
}
