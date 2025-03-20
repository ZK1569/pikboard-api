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

func (self *User) CreateUser(username, email string, password [64]byte) (*model.User, error) {
	user := model.User{
		Username: username,
		Email:    email,
		Password: password[:],
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

func (self *User) GetUserByEmailAndPassword(email string, password [64]byte) (*model.User, error) {
	var user model.User

	result := self.db.DB.First(&user, "email = ? AND password = ?", email, password[:])
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.NotFound
		}
		return nil, result.Error
	}

	return &user, nil
}

func (self *User) UpdateUserSession(user *model.User, token string) error {
	user.Session = &token

	result := self.db.DB.Save(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
