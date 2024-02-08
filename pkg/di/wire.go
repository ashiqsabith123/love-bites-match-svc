//go:build wireinject
// +build wireinject

package di

import (
	auth "github.com/ashiqsabith123/match-svc/pkg/clients/auth_client"
	chat "github.com/ashiqsabith123/match-svc/pkg/clients/chat_client"
	notification "github.com/ashiqsabith123/match-svc/pkg/clients/notification_client"
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
		notification.NewNotificationClient,
		chat.NewChatClient,
		service.NewUserService,
		utils.NewS3Client,
	)

	return service.UserService{}

}
