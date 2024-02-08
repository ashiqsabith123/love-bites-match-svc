package interfaces

import (
	"github.com/ashiqsabith123/love-bytes-proto/match/pb"
	"github.com/ashiqsabith123/match-svc/pkg/helper/responses"
)

type UserUsecase interface {
	SaveAndUploadPhotos(stream pb.MatchService_UplaodPhotosServer) error
	SaveUserPrefrences(req *pb.UserPrefrencesRequest) error
	FindMatches(req *pb.UserIdRequest) (responses.Result, error)
	CreateIntrest(req *pb.IntrestRequest) error
	GetIntrests(req *pb.UserIdRequest) ([]responses.Interests, error)
	ChangeIntrestRequestStatus(req *pb.ChangeInterestRequest) error
}
