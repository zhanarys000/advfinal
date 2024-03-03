package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role,omitempty"`
	Books    []Book `gorm:"many2many:user_books;" json:"books,omitempty"`
}

type Book struct {
	gorm.Model
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Isbn        string `json:"isbn,omitempty"`
	Genre       string `json:"genre,omitempty"`
	Price       string `json:"price,omitempty"`
	Users       []User `gorm:"many2many:user_books;" json:"users,omitempty"`
}

type Subscribers struct {
	Email string `json:"email,omitempty"`
}
type UserBook struct {
	UserID uint
	BookID uint
}
