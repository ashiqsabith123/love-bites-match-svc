package domain

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type UserPhotos struct {
	gorm.Model
	UserID uint           `gorm:"unique;not null"`
	Photos pq.StringArray `gorm:"type:text[]"`
}

type UserPreferences struct {
	gorm.Model
	UserId        uint       `gorm:"unique;not null"`
	UserPhotos    UserPhotos `gorm:"foreignKey:UserId"`
	Height        string     `gorm:"not null"`
	MaritalStatus string     `gorm:"not null"`
	Faith         string     `gorm:"not null"`
	MotherTounge  string     `gorm:"not null"`
	SmokeStatus   string     `gorm:"not null"`
	AlcoholStatus string     `gorm:"not null"`
	SettleStatus  string     `gorm:"not null"`
	Hobbies       string     `gorm:"not null"`
	TeaPerson     string     `gorm:"not null"`
	LoveLanguage  string     `gorm:"not null"`
}

type IntrestRequests struct {
	gorm.Model
	SenderID  uint `gorm:"not null"`
	ReciverID uint `gorm:"not null"`
	Status    string `gorm:"default:'P'"`
}
