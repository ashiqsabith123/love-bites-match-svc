package usecase

import (
	repo "github.com/ashiqsabith123/user-details-svc/pkg/repository/interface"
	interfaces "github.com/ashiqsabith123/user-details-svc/pkg/usecase/interface"
)

type UserUsecase struct {
	UserRepo repo.UserRepo
}

func NewUserUsecase(repo repo.UserRepo) interfaces.UserUsecase {
	return &UserUsecase{UserRepo: repo}
}
