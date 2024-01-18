package repository

import (
	"errors"

	"github.com/ashiqsabith123/match-svc/pkg/domain"
	interfaces "github.com/ashiqsabith123/match-svc/pkg/repository/interface"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type UserRepo struct {
	Postgres *gorm.DB
}

func NewUserRepo(db *gorm.DB) interfaces.UserRepo {
	return &UserRepo{Postgres: db}
}

func (U *UserRepo) SavePhotosID(data domain.UserPhotos) error {

	err := U.Postgres.Create(&data).Error

	if err != nil {
		return err
	}

	return nil

}

func (U *UserRepo) SaveUserPrefrences(data domain.UserPreferences) error {

	err := U.Postgres.Create(&data).Error

	if err != nil {
		return err
	}

	return nil

}

func (U *UserRepo) GetUserPrefrencesByID(ids []int32) (usersPrefrences []domain.UserPreferences, err error) {

	query := "SELECT * FROM user_preferences WHERE user_id = ANY($1)"

	if err := U.Postgres.Raw(query, pq.Array(ids)).Scan(&usersPrefrences).Error; err != nil {
		return []domain.UserPreferences{}, errors.New("error while fetching user prefrences by id: " + err.Error())
	}

	return usersPrefrences, nil
}

func (U *UserRepo) GetUsersPhotosByID(ids []int32) (userPhotos []domain.UserPhotos, err error) {

	query := "SELECT * FROM user_photos WHERE user_id = ANY($1)"

	if err := U.Postgres.Raw(query, pq.Array(ids)).Scan(&userPhotos).Error; err != nil {
		return []domain.UserPhotos{}, errors.New("error while fetching user prefrences by id: " + err.Error())
	}

	return userPhotos, nil
}

func (U *UserRepo) CreateIntrests(intrest domain.IntrestRequests) error {

	err := U.Postgres.Create(&intrest).Error

	if err != nil {
		return err
	}

	return nil
}
