package service

import (
	"context"
	"net/http"
	"sort"

	logs "github.com/ashiqsabith123/love-bytes-proto/log"
	"github.com/ashiqsabith123/love-bytes-proto/match/pb"
	usecase "github.com/ashiqsabith123/match-svc/pkg/usecase/interface"
	"google.golang.org/protobuf/proto"
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

func (U *UserService) GetMatchedUsers(ctx context.Context, req *pb.UserIdRequest) (*pb.MatchResponse, error) {
	matches, err := U.UserUsecase.FindMatches(req)

	if err != nil {
		logs.ErrLog.Println("Error while getting user matches", err)
		return &pb.MatchResponse{
			Code:    500,
			Message: "Server error",
			Error: &anypb.Any{
				Value: []byte(err.Error()),
			},
		}, nil
	}

	data := make([]*pb.MatchedUsers, len(matches.Result))

	for i, v := range matches.Result {
		mathces := &pb.MatchedUsers{
			UserID:     int32(v.UserID),
			Name:       v.Name,
			Age:        int32(v.Age),
			Place:      v.Place,
			MatchScore: v.MatchScore,
			UserImages: v.Photos,
		}

		data[i] = mathces
	}

	matchedUsers := pb.MatchedUsersResponse{
		MatchedUsers: data,
	}

	dataInBytes, err := proto.Marshal(&matchedUsers)
	if err != nil {
		return &pb.MatchResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed while marshaling",
			Error: &anypb.Any{
				Value: []byte(err.Error()),
			},
		}, nil
	}

	return &pb.MatchResponse{
		Code:    http.StatusOK,
		Message: "Data fetched succesfully",
		Data: &anypb.Any{
			Value: dataInBytes,
		},
	}, nil

}

func (U *UserService) CreateIntrests(ctx context.Context, intrest *pb.IntrestRequest) (*pb.MatchResponse, error) {
	err := U.UserUsecase.CreateIntrest(intrest)

	if err != nil {
		logs.ErrLog.Println("Error while create intrest", err)
		return &pb.MatchResponse{
			Code:    500,
			Message: "Server error",
			Error: &anypb.Any{
				Value: []byte(err.Error()),
			},
		}, nil
	}

	return &pb.MatchResponse{
		Code:    http.StatusCreated,
		Message: "Intrest request created succesfully",
	}, nil
}

func (U *UserService) GetAllInteretsRequests(ctx context.Context, req *pb.UserIdRequest) (*pb.MatchResponse, error) {
	intrests, err := U.UserUsecase.GetIntrests(req)

	if err != nil {
		logs.ErrLog.Println("Error while create intrest", err)
		return &pb.MatchResponse{
			Code:    500,
			Message: "Server error",
			Error: &anypb.Any{
				Value: []byte(err.Error()),
			},
		}, nil
	}

	sort.Slice(intrests, func(i, j int) bool {
		return intrests[i].CreatedAt.After(intrests[j].CreatedAt)
	})

	data := make([]*pb.Interest, len(intrests))

	for i, v := range intrests {
		intrests := &pb.Interest{
			
			UserID: uint32(v.UserID),
			Name:   v.Name,
			Photo:  v.Photo,
			Time:   v.CreatedAt.Format("03:04 PM, Mon 02 Jan 2006"),
			Status: v.Status,
		}

		data[i] = intrests
	}

	intrestRequets := pb.IntrestRequests{
		IntrestRequest: data,
	}

	dataInBytes, err := proto.Marshal(&intrestRequets)
	if err != nil {
		return &pb.MatchResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed while marshaling",
			Error: &anypb.Any{
				Value: []byte(err.Error()),
			},
		}, nil
	}

	return &pb.MatchResponse{
		Code:    http.StatusOK,
		Message: "Data fetched succesfully",
		Data: &anypb.Any{
			Value: dataInBytes,
		},
	}, nil

}
