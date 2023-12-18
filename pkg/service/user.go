package service

import (
	"fmt"
	"os"

	"github.com/ashiqsabith123/love-bytes-proto/match/pb"
	usecase "github.com/ashiqsabith123/user-details-svc/pkg/usecase/interface"
)

type UserService struct {
	UserUsecase usecase.UserUsecase
	pb.UnimplementedMatchServiceServer
}

func NewUserService(usecase usecase.UserUsecase) UserService {
	return UserService{UserUsecase: usecase}
}

var i int

func (U *UserService) UplaodPhotos(stream pb.MatchService_UplaodPhotosServer) error {

	var data []byte

	for {

		req, _ := stream.Recv()

		data = append(data, req.ImageData...)

		if req.LastChunk {

			dest, _ := os.Create("image" + fmt.Sprint(i) + ".jpeg")
			i++

			_, err := dest.Write(data)
			if err != nil {
				fmt.Println("errrrrrr", err)
			}

			err = dest.Close()
			if err != nil {
				fmt.Println("err", err)
			}

		}

	}

	// dest, _ := os.Create("image.jpeg")

	// for {
	// 	req, _ := stream.Recv()

	// 	fmt.Println(req.ImageData, req.LastChunk)

	// 	if !req.LastChunk {
	// 		_, err := dest.Write(req.ImageData)
	// 		if err != nil {
	// 			fmt.Println("errrrrrr", err)
	// 		}
	// 	} else {
	// 		_, err := dest.Write(req.ImageData)
	// 		if err != nil {
	// 			fmt.Println("errrrrrr", err)
	// 		}
	// 		dest, _ = os.Create("image1.jpeg")
	// 	}

	// }

	return nil

}
