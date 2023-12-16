package models

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ricardoalcantara/web-notify/internal/utils"
	"github.com/rs/zerolog/log"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func ConnectDataBase() {
	dbUrl := os.Getenv("DB_URL")
	envDialector := os.Getenv("DB_DIALECTOR")
	var err error
	var dialector gorm.Dialector
	switch strings.ToLower(envDialector) {
	case "sqlite":
		dialector = sqlite.Open(dbUrl)
	case "mysql":
		dialector = mysql.Open(dbUrl)
	case "postgres":
		dialector = postgres.Open(dbUrl)
	default:
		log.Fatal().Err(err).Msg("connection error:")
	}
	db, err = gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("connection error:")
	} else {
		log.Debug().Msg("Db Connected")
	}

	migrate()
	createAdmin()
}

func migrate() {
	db.AutoMigrate(&Client{})
	db.AutoMigrate(&Message{})
}

func createAdmin() {
	if err := db.Take(&Client{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		clientId := utils.GetEnvOr("CLIENT_ID", func() string {
			return utils.GenString(50)
		})
		clientSecret := utils.GetEnvOr("CLIENT_SECRET", func() string {
			return utils.GenString(100)
		})

		log.Debug().Msg("Admin Created")

		client := Client{
			Name:         "Admin",
			ClientId:     clientId,
			ClientSecret: clientSecret,
		}

		client.Save()

		fmt.Printf("ClientId: %s\n", client.ClientId)
		fmt.Printf("ClientSecret: %s\n", client.ClientSecret)
	}
}
