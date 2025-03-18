package util

import (
	"fmt"
	"log"
	"sync"

	errs "github.com/zk1569/pikboard-api/src/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database interface {
	DB() any
}

type DatabasePostgres struct {
	db *gorm.DB
}

var lock = &sync.Mutex{}
var singleInstance *DatabasePostgres

func GetDatabasePostgresInstance() Database {
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
			db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

			if err != nil {
				log.Fatalf("%s - ❌ An error occurred when connecting to the database: %s", errs.ErrDatabase, err)
			}

			singleInstance = &DatabasePostgres{
				db: db,
			}
		}
	}

	return singleInstance
}

func (self *DatabasePostgres) DB() any {
	return self.db
}
