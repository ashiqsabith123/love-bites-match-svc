package repository

import (
	"errors"

	"github.com/ashiqsabith123/match-svc/pkg/domain"
	"github.com/ashiqsabith123/match-svc/pkg/helper/responses"
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

func (U *UserRepo) CreateIntrestsAndReturnID(intrest domain.IntrestRequests) (int, error) {

	err := U.Postgres.Create(&intrest).Error

	if err != nil {
		return 0, err
	}

	return int(intrest.ID), nil
}

func (U *UserRepo) GetIntrestRequestAndPhotoById(id uint) (userIntrests []responses.Interests, err error) {

	query :=
		"SELECT users.id,users.user_id AS user_id, users.created_at, users.status, (user_photos.photos)[1] AS photo FROM (SELECT id, receiver_id AS user_id, created_at, status FROM intrest_requests WHERE sender_id = $1 AND status = 'A' UNION SELECT id, sender_id AS user_id, created_at, status FROM intrest_requests WHERE receiver_id = $2) AS users INNER JOIN user_photos ON user_photos.user_id = users.user_id;"
	if err := U.Postgres.Raw(query, id, id).Scan(&userIntrests).Error; err != nil {
		return userIntrests, err
	}

	return userIntrests, nil
}

func (U *UserRepo) GetUserPhotoByID(id int) (photo string, err error) {

	query := "SELECT photos[1] FROM user_photos WHERE user_id = $1"

	if err := U.Postgres.Raw(query, id).Scan(&photo).Error; err != nil {
		return "", err

	}

	return photo, nil
}
