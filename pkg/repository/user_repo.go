package repository

import (
	"github.com/ashiqsabith123/user-details-svc/pkg/domain"
	interfaces "github.com/ashiqsabith123/user-details-svc/pkg/repository/interface"
	"gorm.io/gorm"
)

type UserRepo struct {
	Postgres *gorm.DB
}

func NewUserRepo(db *gorm.DB) interfaces.UserRepo {
	return &UserRepo{Postgres: db}
}

func (U *UserRepo) SavePhotos(data domain.UserPhotos) error {

	err := U.Postgres.Create(&data).Error

	if err != nil {
		return err
	}

	return nil

}
