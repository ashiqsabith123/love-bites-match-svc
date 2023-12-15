package service

import (
	usecase "github.com/ashiqsabith123/user-details-svc/pkg/usecase/interface"
)

type UserService struct {
	UserUsecase usecase.UserUsecase
}

func NewUserService(usecase usecase.UserUsecase) UserService {
	return UserService{UserUsecase: usecase}
}
