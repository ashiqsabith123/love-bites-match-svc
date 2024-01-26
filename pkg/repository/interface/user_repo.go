package interfaces

import (
	"github.com/ashiqsabith123/match-svc/pkg/domain"
	"github.com/ashiqsabith123/match-svc/pkg/helper/responses"
)

type UserRepo interface {
	SavePhotosID(data domain.UserPhotos) error
	SaveUserPrefrences(data domain.UserPreferences) error
	GetUserPrefrencesByID(ids []int32) (usersPrefrences []domain.UserPreferences, err error)
	GetUsersPhotosByID(ids []int32) (userPhotos []domain.UserPhotos, err error)
	CreateIntrests(intrest domain.IntrestRequests) error
	GetIntrestRequestAndPhotoById(id uint) (userIntrests []responses.Interests, err error)
	GetUserPhotoByID(id int) (photo string, err error)
}
