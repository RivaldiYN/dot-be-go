package entity

import (
	"time"

	"gorm.io/gorm"
)

// Book represents the book model
type Book struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"size:255;not null"`
	Author      string         `json:"author" gorm:"size:100;not null"`
	ISBN        string         `json:"isbn" gorm:"size:20;uniqueIndex"`
	PublishYear int            `json:"publish_year"`
	Description string         `json:"description" gorm:"type:text"`
	Categories  []Category     `json:"categories,omitempty" gorm:"many2many:book_categories;"`
	UserID      uint           `json:"user_id" gorm:"not null"`
	User        User           `json:"-" gorm:"foreignKey:UserID"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name for Book
func (Book) TableName() string {
	return "books"
}
