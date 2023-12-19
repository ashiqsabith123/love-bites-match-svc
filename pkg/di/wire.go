//go:build wireinject
// +build wireinject

package di

import (
	"github.com/ashiqsabith123/user-details-svc/pkg/config"
	"github.com/ashiqsabith123/user-details-svc/pkg/db"
	"github.com/ashiqsabith123/user-details-svc/pkg/repository"
	"github.com/ashiqsabith123/user-details-svc/pkg/service"
	"github.com/ashiqsabith123/user-details-svc/pkg/usecase"
	"github.com/ashiqsabith123/user-details-svc/pkg/utils"
	"github.com/google/wire"
)

func IntializeService(config config.Config) service.UserService {

	wire.Build(
		db.ConnectToDatabase,
		repository.NewUserRepo,
		usecase.NewUserUsecase,
		service.NewUserService,
		utils.NewS3Client,
	)

	return service.UserService{}

}
