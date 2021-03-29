package database

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;unique;not null" json:"username"`
	Email    string `gorm:"uniqueIndex;unique;not null" json:"email"`
	Password string `gorm:"not null" json:"password"`
}
