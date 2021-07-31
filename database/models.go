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

type HTTPMethod struct {
	gorm.Model
	Method string `gorm:"uniqueIndex;unique;not null" json:"method"`
}

type Rule struct {
	gorm.Model     `yaml:"-"`
	Path           string       `gorm:"index;not null" json:"path" yaml:"Path"`
	MetaKey        string       `json:"meta_key" yaml:"MetaKey"`
	PathPrefix     bool         `gorm:"default:false" json:"path_prefix" yaml:"path_prefix"`
	MetaLocationID uint         `json:"meta_location_id" yaml:"MetaLocationID"`
	MetaLocation   MetaLocation `gorm:"constraint:OnDelete:CASCADE;" yaml:"-"`
	HTTPMethodID   uint         `json:"http_method_id" yaml:"HTTPMethodID"`
	HTTPMethod     HTTPMethod   `gorm:"constraint:OnDelete:CASCADE;" yaml:"-"`
}

type AssignedRules struct {
	gorm.Model
	RuleID    uint   `json:"rule_id"`
	Rule      Rule   `gorm:"constraint:OnDelete:CASCADE;"`
	UserID    uint   `json:"user_id"`
	User      User   `gorm:"constraint:OnDelete:CASCADE;"`
	MetaValue string `json:"meta_value"`
}
