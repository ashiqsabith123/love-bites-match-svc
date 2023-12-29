package service

import (
	"context"

	logs "github.com/ashiqsabith123/love-bytes-proto/log"
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
		logs.ErrLog.Println("Error while uploading photo", err)
		return stream.SendAndClose(&pb.MatchResponse{
			Code:    500,
			Message: "Error while uploading photos",
			Error: &anypb.Any{
				Value: []byte(err.Error()),
			},
		})
	}

	logs.GenLog.Println("Photo upload succecsfully")

	return stream.SendAndClose(&pb.MatchResponse{
		Code:    201,
		Message: "photo upload succecsfully",
	})

}

func (U *UserService) SaveUserPrefrences(ctx context.Context, req *pb.UserPrefrencesRequest) (*pb.MatchResponse, error) {

	err := U.UserUsecase.SaveUserPrefrences(req)

	if err != nil {
		logs.ErrLog.Println("Error while saving user prefrences", err)
		return &pb.MatchResponse{
			Code:    500,
			Message: "Server error",
			Error: &anypb.Any{
				Value: []byte(err.Error()),
			},
		}, nil
	}

	logs.GenLog.Println("User prefrences added succecsfully")
	return &pb.MatchResponse{
		Code:    201,
		Message: "User prefrences added succecsfully",
	}, nil
}
