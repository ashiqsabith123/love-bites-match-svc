//go:build wireinject
// +build wireinject

package di

import (
	"github.com/ashiqsabith123/match-svc/pkg/clients/auth"
	"github.com/ashiqsabith123/match-svc/pkg/config"
	"github.com/ashiqsabith123/match-svc/pkg/db"
	"github.com/ashiqsabith123/match-svc/pkg/repository"
	"github.com/ashiqsabith123/match-svc/pkg/service"
	"github.com/ashiqsabith123/match-svc/pkg/usecase"
	"github.com/ashiqsabith123/match-svc/pkg/utils"
	"github.com/google/wire"
)

func IntializeService(config config.Config) service.UserService {

	wire.Build(
		db.ConnectToDatabase,
		repository.NewUserRepo,
		usecase.NewUserUsecase,
		auth.NewAuthClient,
		service.NewUserService,
		utils.NewS3Client,
	)

	return service.UserService{}

}
