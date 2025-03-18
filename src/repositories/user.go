package repository

import (
	"errors"

	errs "github.com/zk1569/pikboard-api/src/errors"
	model "github.com/zk1569/pikboard-api/src/models"
	util "github.com/zk1569/pikboard-api/src/utils"
	"gorm.io/gorm"
)

type User struct {
	db *util.DatabasePostgres
}

var singleUserInstance UserInterface

func GetUserInstance() UserInterface {
	if singleUserInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleUserInstance == nil {
			singleUserInstance = &User{
				db: util.GetDatabasePostgresInstance(),
			}
		}
	}

	return singleUserInstance
}

func (self *User) CreateUser(username, email, password string) (*model.User, error) {
	user := model.User{
		Username: username,
		Email:    email,
		Password: password,
	}

	result := self.db.DB.Create(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, errs.AlreadyExists
		}
		return nil, result.Error
	}

	return &user, nil
}
