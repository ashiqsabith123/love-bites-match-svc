package interfaces

import "github.com/ashiqsabith123/match-svc/pkg/domain"

type UserRepo interface {
	SavePhotosID(data domain.UserPhotos) error
	SaveUserPrefrences(data domain.UserPreferences) error
}
