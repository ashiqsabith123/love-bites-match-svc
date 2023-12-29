package interfaces

import "github.com/ashiqsabith123/love-bytes-proto/match/pb"

type UserUsecase interface {
	SaveAndUploadPhotos(stream pb.MatchService_UplaodPhotosServer) error
	SaveUserPrefrences(req *pb.UserPrefrencesRequest) error
}
