package usecase

import (
	"io"
	"sync"

	"github.com/ashiqsabith123/love-bytes-proto/match/pb"
	"github.com/ashiqsabith123/user-details-svc/pkg/domain"
	logs "github.com/ashiqsabith123/love-bytes-proto/log"
	repo "github.com/ashiqsabith123/user-details-svc/pkg/repository/interface"
	interfaces "github.com/ashiqsabith123/user-details-svc/pkg/usecase/interface"
	utils "github.com/ashiqsabith123/user-details-svc/pkg/utils/interface"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type UserUsecase struct {
	UserRepo repo.UserRepo
	Utils    utils.Utils
}

func NewUserUsecase(repo repo.UserRepo, utils utils.Utils) interfaces.UserUsecase {
	return &UserUsecase{UserRepo: repo, Utils: utils}
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
