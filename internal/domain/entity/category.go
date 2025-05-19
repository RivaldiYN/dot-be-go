package entity

import (
	"time"

	"gorm.io/gorm"
)

// Category represents the book category model
type Category struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:100;not null;uniqueIndex"`
	Description string         `json:"description" gorm:"size:255"`
	Books       []Book         `json:"books,omitempty" gorm:"many2many:book_categories;"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name for Category
func (Category) TableName() string {
	return "categories"
}
