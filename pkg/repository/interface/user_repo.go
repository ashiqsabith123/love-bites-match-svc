package interfaces

import "github.com/ashiqsabith123/user-details-svc/pkg/domain"

type UserRepo interface {
	SavePhotosID(data domain.UserPhotos) error
	SaveUserPrefrences(data domain.UserPreferences) error
}
