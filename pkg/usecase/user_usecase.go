package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	authPb "github.com/ashiqsabith123/love-bytes-proto/auth/pb"
	chatPb "github.com/ashiqsabith123/love-bytes-proto/chat/pb"
	logs "github.com/ashiqsabith123/love-bytes-proto/log"
	"github.com/ashiqsabith123/love-bytes-proto/match/pb"
	notifiPb "github.com/ashiqsabith123/love-bytes-proto/notifications/pb"
	authClient "github.com/ashiqsabith123/match-svc/pkg/clients/auth_client/interface"
	chatClient "github.com/ashiqsabith123/match-svc/pkg/clients/chat_client/interface"
	notifiClient "github.com/ashiqsabith123/match-svc/pkg/clients/notification_client/interface"
	"github.com/ashiqsabith123/match-svc/pkg/domain"
	"github.com/ashiqsabith123/match-svc/pkg/helper/responses"
	repo "github.com/ashiqsabith123/match-svc/pkg/repository/interface"
	interfaces "github.com/ashiqsabith123/match-svc/pkg/usecase/interface"
	utils "github.com/ashiqsabith123/match-svc/pkg/utils/interface"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/proto"
)

type UserUsecase struct {
	UserRepo           repo.UserRepo
	Utils              utils.Utils
	AuthClient         authPb.AuthServiceClient
	NotificationClient notifiPb.NotificationServiceClient
	ChatClient         chatPb.ChatServiceClient
}

func NewUserUsecase(repo repo.UserRepo, utils utils.Utils, authClient authClient.AuthClient, notifiClient notifiClient.NotificationClient, chatClient chatClient.ChatClient) interfaces.UserUsecase {
	return &UserUsecase{UserRepo: repo, Utils: utils, AuthClient: authClient.GetClient(), NotificationClient: notifiClient.GetClient(), ChatClient: chatClient.GetClient()}
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

func (U *UserUsecase) FindMatches(req *pb.UserIdRequest) (responses.Result, error) {

	resp, _ := U.AuthClient.GetUserByID(context.TODO(), &authPb.UserIDRequest{UserID: req.UserID})

	if resp.Data == nil {
		return responses.Result{}, errors.New("coudnt fetch user data by id")
	}

	var person1 authPb.UserRepsonse

	if resp.Data != nil {
		if err := proto.Unmarshal(resp.Data.Value, &person1); err != nil {
			return responses.Result{}, err
		}
	}

	gender := "M"

	if person1.Gender == "M" {
		gender = "F"
	}

	resp, err := U.AuthClient.GetUsersByGender(context.TODO(), &authPb.UserGenderRequest{Gender: gender})

	if err != nil {
		return responses.Result{}, err
	}

	var person2sDataByGender authPb.UserResponses

	if resp.Data != nil {
		if err := proto.Unmarshal(resp.Data.Value, &person2sDataByGender); err != nil {
			return responses.Result{}, err
		}
	}

	userIDs := []int32{}

	for _, v := range person2sDataByGender.UserRepsonses {
		userIDs = append(userIDs, v.UserID)
	}

	person2sPrefrences, err := U.UserRepo.GetUserPrefrencesByID(userIDs)
	if err != nil {
		return responses.Result{}, err
	}

	person1Prefrences, err := U.UserRepo.GetUserPrefrencesByID([]int32{req.UserID})
	if err != nil {
		return responses.Result{}, err
	}

	match, err := U.Utils.MakeMatchesByPrefrences(&person1, person2sDataByGender.UserRepsonses, person1Prefrences, person2sPrefrences)

	if err != nil {
		return responses.Result{}, err
	}

	userIDs = []int32{}

	for _, v := range match.Result {
		userIDs = append(userIDs, int32(v.UserID))
	}

	userPhotos, err := U.UserRepo.GetUsersPhotosByID(userIDs)
	if err != nil {
		return responses.Result{}, err
	}

	for i, v := range match.Result {
		for _, k := range userPhotos {
			if v.UserID == k.UserID {
				photo := make([]*pb.Images, len(k.Photos))

				for i := 0; i < len(k.Photos); i++ {
					photo[i] = &pb.Images{ImageId: k.Photos[i]}
				}
				match.Result[i].Photos = photo
				break
			}
		}

	}

	return match, nil

}

func (U *UserUsecase) CreateIntrest(req *pb.IntrestRequest) error {
	var intrestReq domain.IntrestRequests
	copier.Copy(&intrestReq, &req)

	interestID, err := U.UserRepo.CreateIntrestsAndReturnID(intrestReq)
	if err != nil {
		return err
	}

	// wg := sync.WaitGroup{}
	// errCh := make(chan error, 2)
	// name := make(chan string, 1)  // Initialize with buffer size 1
	// photo := make(chan string, 1) // Initialize with buffer size 1

	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()

	// 	resp, err := U.AuthClient.GetUserByID(context.Background(), &authPb.UserIDRequest{UserID: int32(req.SenderID)})
	// 	if err != nil || resp == nil || resp.Data == nil {
	// 		errCh <- errors.New("data not found from the server")
	// 		return
	// 	}

	// 	var userDetails authPb.UserRepsonse
	// 	if err := proto.Unmarshal(resp.Data.Value, &userDetails); err != nil {
	// 		errCh <- errors.New("Data not found from the server" + err.Error())
	// 		return
	// 	}

	// 	name <- userDetails.Fullname
	// }()

	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()

	// 	pto, err := U.UserRepo.GetUserPhotoByID(int(req.SenderID))
	// 	if err != nil {
	// 		errCh <- err
	// 		return
	// 	}

	// 	photo <- pto
	// }()

	// wg.Wait()
	// close(errCh)

	// Retrieve values from channels
	// senderName := <-name
	// senderPhoto := <-photo

	// fmt.Println(senderName, senderPhoto)

	// for err := range errCh {
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	senderName, senderPhoto, err := U.GetNameAndPhotoOfUser(int(req.SenderID))

	if err != nil {
		return err
	}

	resp, err := U.NotificationClient.CreateNotification(context.Background(), &notifiPb.Notification{

		SenderID:   req.SenderID,
		ReceiverID: req.ReceiverID,
		Name:       senderName,
		Image:      senderPhoto,
		Type:       "IR",
		Status:     "P",
		CommonID:   uint32(interestID),
	})
	if err != nil || resp.Err != "" {
		return errors.New("Failed to create notification: " + err.Error())
	}

	return nil
}

func (U *UserUsecase) GetIntrests(req *pb.UserIdRequest) ([]responses.Interests, error) {

	wg := sync.WaitGroup{}

	data, err := U.UserRepo.GetIntrestRequestAndPhotoById(uint(req.UserID))

	if err != nil {
		return []responses.Interests{}, err
	}

	ch := make(chan error, len(data))

	for i, v := range data {

		wg.Add(1)

		go func(index int, id int32, wg *sync.WaitGroup, ch chan error) {
			resp, _ := U.AuthClient.GetUserByID(context.Background(), &authPb.UserIDRequest{UserID: id})

			if resp == nil || resp.Data == nil {

				ch <- errors.New("data not found from the server")
				return

			}

			var userDetails authPb.UserRepsonse

			if err := proto.Unmarshal(resp.Data.Value, &userDetails); err != nil {
				ch <- errors.New("Data not found from the server" + err.Error())
				return
			}

			data[index].Name = userDetails.Fullname

			wg.Done()

		}(i, int32(v.UserID), &wg, ch)
	}

	wg.Wait()
	close(ch)

	for err := range ch {
		if err != nil {
			return []responses.Interests{}, err
		}
	}

	return data, nil

}

func (U *UserUsecase) ChangeIntrestRequestStatus(req *pb.ChangeInterestRequest) error {

	interest, err := U.UserRepo.ChangeIntrestRequestStatus(int(req.Id), req.Status)

	if err != nil {
		return err
	}

	if req.Status == "S" {

		resp, _ := U.ChatClient.CreateChatRoom(context.Background(), &chatPb.ChatRoomRequest{
			UserID1: uint32(interest.SenderID),
			UserID2: uint32(interest.ReceiverID),
		})

		fmt.Println(resp)

		if resp != nil && resp.Error != nil {
			return errors.New(string(resp.Error.GetValue()))
		}
	}

	name, photo, err := U.GetNameAndPhotoOfUser(int(interest.ReceiverID))

	if err != nil {
		return err
	}

	resp, err := U.NotificationClient.CreateNotification(context.Background(), &notifiPb.Notification{
		ReceiverID: uint32(interest.SenderID),
		SenderID:   uint32(interest.ReceiverID),
		Name:       name,
		Image:      photo,
		Type:       "IR",
		Status:     req.Status,
		CommonID:   uint32(interest.ID),
	})

	if err != nil || resp.Err != "" {
		return errors.New("Failed to create notification: " + err.Error())
	}

	return nil
}

func (U *UserUsecase) GetNameAndPhotoOfUser(id int) (string, string, error) {

	wg := sync.WaitGroup{}
	errCh := make(chan error, 2)
	name := make(chan string, 1)  // Initialize with buffer size 1
	photo := make(chan string, 1) // Initialize with buffer size 1

	wg.Add(1)
	go func() {
		defer wg.Done()

		resp, err := U.AuthClient.GetUserByID(context.Background(), &authPb.UserIDRequest{UserID: int32(id)})
		if err != nil || resp == nil || resp.Data == nil {
			errCh <- errors.New("data not found from the server")
			return
		}

		var userDetails authPb.UserRepsonse
		if err := proto.Unmarshal(resp.Data.Value, &userDetails); err != nil {
			errCh <- errors.New("Data not found from the server" + err.Error())
			return
		}

		name <- userDetails.Fullname
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		pto, err := U.UserRepo.GetUserPhotoByID(int(id))
		if err != nil {
			errCh <- err
			return
		}

		photo <- pto
	}()

	wg.Wait()
	close(errCh)

	// Retrieve values from channels
	senderName := <-name
	senderPhoto := <-photo

	fmt.Println(senderName, senderPhoto)

	for err := range errCh {
		if err != nil {
			return "", "", err
		}
	}

	return senderName, senderPhoto, nil

}
