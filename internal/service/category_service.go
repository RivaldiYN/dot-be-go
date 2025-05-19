package service

import (
	"dot-be-go/internal/domain/entity"
	"dot-be-go/internal/domain/repository"
)

// CategoryRequest represents category request data
type CategoryRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description" validate:"max=255"`
}

// CategoryService handles category operations
type CategoryService interface {
	Create(req *CategoryRequest) (*entity.Category, error)
	GetAll() ([]entity.Category, error)
	GetByID(id uint) (*entity.Category, error)
	Update(id uint, req *CategoryRequest) (*entity.Category, error)
	Delete(id uint) error
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

// NewCategoryService creates a new category service
func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

// Create creates a new category
func (s *categoryService) Create(req *CategoryRequest) (*entity.Category, error) {
	category := &entity.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

// GetAll returns all categories
func (s *categoryService) GetAll() ([]entity.Category, error) {
	return s.categoryRepo.FindAll()
}

// GetByID returns a category by ID
func (s *categoryService) GetByID(id uint) (*entity.Category, error) {
	return s.categoryRepo.FindByID(id)
}

// Update updates a category
func (s *categoryService) Update(id uint, req *CategoryRequest) (*entity.Category, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	category.Name = req.Name
	category.Description = req.Description

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

// Delete deletes a category
func (s *categoryService) Delete(id uint) error {
	return s.categoryRepo.Delete(id)
}
