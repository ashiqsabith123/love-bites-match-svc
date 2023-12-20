package usecase

import (
	"io"
	"sync"

	"github.com/ashiqsabith123/love-bytes-proto/match/pb"
	"github.com/ashiqsabith123/user-details-svc/pkg/domain"
	repo "github.com/ashiqsabith123/user-details-svc/pkg/repository/interface"
	interfaces "github.com/ashiqsabith123/user-details-svc/pkg/usecase/interface"
	utils "github.com/ashiqsabith123/user-details-svc/pkg/utils/interface"
	"github.com/google/uuid"
)

type UserUsecase struct {
	UserRepo repo.UserRepo
	Utils    utils.Utils
}

func NewUserUsecase(repo repo.UserRepo, utils utils.Utils) interfaces.UserUsecase {
	return &UserUsecase{UserRepo: repo, Utils: utils}
}

var i int

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

			// dest, _ := os.Create("image" + fmt.Sprint(i) + ".jpeg")
			// i++

			// _, err := dest.Write(data)
			// if err != nil {
			// 	fmt.Println("errrrrrr", err)
			// }

			// err = dest.Close()
			// if err != nil {
			// 	fmt.Println("err", err)
			// }

			id := uuid.New()
			imageId := id.String()

			// wg.Add(1)
			// go func() {
			// 	defer wg.Done()
			err = U.Utils.UploadPhotos(imageId, data)
			if err != nil {
				return err
			}
			// 	if err != nil {
			// 		ch <- err
			// 	}
			// }()

			photos.Photos = append(photos.Photos, imageId)
			photos.UserID = req.UserID

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

	err := U.UserRepo.SavePhotos(photos)

	if err != nil {
		return err
	}

	return nil

}
