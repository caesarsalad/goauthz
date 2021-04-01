package database

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;unique;not null" json:"username"`
	Email    string `gorm:"uniqueIndex;unique;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
}

type MetaLocation struct {
	gorm.Model
	MetaLocation string `gorm:"uniqueIndex;unique;not null" json:"meta_location"`
}

type Rule struct {
	gorm.Model
	Path           string       `gorm:"uniqueIndex;unique;not null" json:"path"`
	MetaKey        string       `json:"meta_key"`
	MetaLocationID uint         `json:"meta_location_id"`
	MetaLocation   MetaLocation `gorm:"constraint:OnDelete:CASCADE;"`
}

type AssignedRules struct {
	gorm.Model
	RuleID    uint   `json:"rule_id"`
	Rule      Rule   `gorm:"constraint:OnDelete:CASCADE;"`
	UserID    uint   `json:"user_id"`
	User      User   `gorm:"constraint:OnDelete:CASCADE;"`
	MetaValue string `json:"meta_value"`
}
