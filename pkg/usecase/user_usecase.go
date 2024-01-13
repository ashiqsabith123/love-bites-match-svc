package usecase

import (
	"context"
	"fmt"
	"io"
	"sync"

	authPb "github.com/ashiqsabith123/love-bytes-proto/auth/pb"
	logs "github.com/ashiqsabith123/love-bytes-proto/log"
	"github.com/ashiqsabith123/love-bytes-proto/match/pb"
	authClient "github.com/ashiqsabith123/match-svc/pkg/clients/auth/interface"
	"github.com/ashiqsabith123/match-svc/pkg/domain"
	repo "github.com/ashiqsabith123/match-svc/pkg/repository/interface"
	interfaces "github.com/ashiqsabith123/match-svc/pkg/usecase/interface"
	utils "github.com/ashiqsabith123/match-svc/pkg/utils/interface"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/proto"
)

type UserUsecase struct {
	UserRepo repo.UserRepo
	Utils    utils.Utils
	Client   authPb.AuthServiceClient
}

func NewUserUsecase(repo repo.UserRepo, utils utils.Utils, client authClient.AuthClient) interfaces.UserUsecase {
	return &UserUsecase{UserRepo: repo, Utils: utils, Client: client.GetClient()}
}

func (U *UserUsecase) SaveAndUploadPhotos(stream pb.MatchService_UplaodPhotosServer) error {

	var data []byte
	var photos domain.UserPhotos

	wg := sync.WaitGroup{}
	ch := make(chan error, 4)

	for {

		req, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		data = append(data, req.ImageData...)

		if req.LastChunk {

			id := uuid.New()
			imageId := id.String()

			wg.Add(1)
			go func(imageID string, imageData []byte) {
				defer func() {
					if r := recover(); r != nil {
						logs.ErrLog.Println("Panic occured while uploading image : ", r)
					}
				}()
				defer wg.Done()
				err := U.Utils.UploadPhotos(imageID+".jpeg", imageData)
				if err != nil {
					ch <- err
				}
			}(imageId, data)

			photos.Photos = append(photos.Photos, imageId)
			photos.UserID = uint(req.UserID)

			data = nil

		}

	}

	wg.Wait()
	close(ch)

	for err := range ch {
		if err != nil {
			return err
		}
	}

	err := U.UserRepo.SavePhotosID(photos)

	if err != nil {
		return err
	}

	return nil

}

func (U *UserUsecase) SaveUserPrefrences(req *pb.UserPrefrencesRequest) error {

	var userPreferences domain.UserPreferences
	err := copier.Copy(&userPreferences, req)
	if err != nil {
		return err
	}

	err = U.UserRepo.SaveUserPrefrences(userPreferences)
	if err != nil {
		return err
	}

	return nil
}

func (U *UserUsecase) FindMatches() error {

	resp, err := U.Client.GetUserByID(context.TODO(), &authPb.UserIDRequest{UserID: 2})

	if err != nil {
		return err
	}

	var userData authPb.UserRepsonse

	if resp.Data != nil {
		if err := proto.Unmarshal(resp.Data.Value, &userData); err != nil {
			return err
		}
	}

	gender := "M"

	if userData.Gender == "M" {
		gender = "F"
	}


	resp, err = U.Client.GetUsersByGender(context.TODO(), &authPb.UserGenderRequest{Gender: gender})

	if err != nil {
		return err
	}

	var userDataByGender authPb.UserResponses

	if resp.Data != nil {
		if err := proto.Unmarshal(resp.Data.Value, &userDataByGender); err != nil {
			return err
		}
	}

	fmt.Println(userDataByGender.UserRepsonses)

	return nil

}
