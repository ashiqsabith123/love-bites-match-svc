package interfaces

import "github.com/ashiqsabith123/user-details-svc/pkg/domain"

type UserRepo interface {
	SavePhotos(data domain.UserPhotos) error
}
