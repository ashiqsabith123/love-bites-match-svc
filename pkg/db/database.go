package db

import (
	"fmt"

	logs "github.com/ashiqsabith123/love-bytes-proto/log"
	"github.com/ashiqsabith123/match-svc/pkg/config"
	"github.com/ashiqsabith123/match-svc/pkg/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDatabase(config config.Config) *gorm.DB {
	connstr := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", config.Postgres.Host, config.Postgres.User, config.Postgres.Database, config.Postgres.Port, config.Postgres.Paswword)
	db, err := gorm.Open(postgres.Open(connstr), &gorm.Config{})

	if err != nil {
		logs.ErrLog.Fatal("Failed to connect database - ", err)
		return nil
	}

	db.AutoMigrate(
		domain.UserPhotos{},
		domain.UserPreferences{},
		domain.IntrestRequests{},
	)

	logs.GenLog.Println("Database connected succesfully....")

	return db
}
