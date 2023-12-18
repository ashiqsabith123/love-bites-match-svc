package db

import (
	"fmt"
	"log"

	"github.com/ashiqsabith123/user-details-svc/pkg/config"
	"github.com/ashiqsabith123/user-details-svc/pkg/helper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDatabase(config config.Config) *gorm.DB {
	connstr := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", config.Postgres.Host, config.Postgres.User, config.Postgres.Database, config.Postgres.Port, config.Postgres.Paswword)
	db, err := gorm.Open(postgres.Open(connstr), &gorm.Config{})

	if err != nil {
		log.Fatal(helper.Red("Failed to connect database - ", err))
		return nil
	}

	db.AutoMigrate(
	//domain.UserPhotos{},
	)

	fmt.Println(helper.Green("Database connected succesfully...."))

	return db
}
