package service

import (
	"github.com/ashiqsabith123/love-bytes-proto/match/pb"
	usecase "github.com/ashiqsabith123/user-details-svc/pkg/usecase/interface"
	"google.golang.org/protobuf/types/known/anypb"
)

type UserService struct {
	UserUsecase usecase.UserUsecase
	pb.UnimplementedMatchServiceServer
}

func NewUserService(usecase usecase.UserUsecase) UserService {
	return UserService{UserUsecase: usecase}
}

func (U *UserService) UplaodPhotos(stream pb.MatchService_UplaodPhotosServer) error {

	err := U.UserUsecase.SaveAndUploadPhotos(stream)

	if err != nil {
		return stream.SendAndClose(&pb.MatchResponse{
			Code:    500,
			Message: "Problem while uploading photos",
			Error: &anypb.Any{
				Value: []byte(err.Error()),
			},
		})
	}

	return stream.SendAndClose(&pb.MatchResponse{
		Code:    200,
		Message: "photo upload succecsfully",
	})

}
