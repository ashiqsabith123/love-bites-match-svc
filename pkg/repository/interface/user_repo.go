package interfaces

import "github.com/ashiqsabith123/match-svc/pkg/domain"

type UserRepo interface {
	SavePhotosID(data domain.UserPhotos) error
	SaveUserPrefrences(data domain.UserPreferences) error
	GetUserPrefrencesByID(ids []int32) (usersPrefrences []domain.UserPreferences, err error)
	GetUsersPhotosByID(ids []int32) (userPhotos []domain.UserPhotos, err error)
	CreateIntrests(intrest domain.IntrestRequests) error
}
