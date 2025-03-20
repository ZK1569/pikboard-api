package util

import (
	"fmt"
	"log"
	"sync"

	errs "github.com/zk1569/pikboard-api/src/errors"
	model "github.com/zk1569/pikboard-api/src/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabasePostgres struct {
	DB *gorm.DB
}

var lock = &sync.Mutex{}
var singleInstance *DatabasePostgres

func GetDatabasePostgresInstance() *DatabasePostgres {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()

		if singleInstance == nil {
			dsn := fmt.Sprintf(
				"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
				EnvVariable.Database.Host,
				EnvVariable.Database.Username,
				EnvVariable.Database.Password,
				EnvVariable.Database.Name,
				EnvVariable.Database.Port,
			)
			db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})

			if err != nil {
				log.Fatalf("%s - ❌ An error occurred when connecting to the database: %s", errs.Database, err)
			}

			err = migrate_database(db)
			if err != nil {
				log.Fatalf("%s - ❌ An error occureed during migrations: %s", errs.Database, err)
			}

			singleInstance = &DatabasePostgres{
				DB: db,
			}
		}
	}

	return singleInstance
}

func migrate_database(db *gorm.DB) error {
	log.Printf(" ⚙️ Start migrations ...")
	err := db.AutoMigrate(&model.User{})

	if err != nil {
		return err
	}

	err = db.AutoMigrate(&model.FriendRequest{})
	if err != nil {
		return err
	}

	return nil
}
