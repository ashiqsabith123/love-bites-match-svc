package intrefaces

import (
	authPb "github.com/ashiqsabith123/love-bytes-proto/auth/pb"
	"github.com/ashiqsabith123/match-svc/pkg/domain"
	"github.com/ashiqsabith123/match-svc/pkg/helper/responses"
)

type Utils interface {
	UploadPhotos(key string, image []byte) error
	FilterWithDistance(user *authPb.UserRepsonse, matchUsersList []authPb.UserRepsonse)
	MakeMatchesByPrefrences(person1Data *authPb.UserRepsonse, person2Data []*authPb.UserRepsonse, person1Prefrences []domain.UserPreferences, person2sPrefrences []domain.UserPreferences) (responses.Result, error)
}
