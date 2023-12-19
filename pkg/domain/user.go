package domain

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)



type UserPhotos struct {
	gorm.Model
	UserID int32          `gorm:"unique;not null" `
	Photos pq.StringArray `gorm:"type:text[]"`
}
