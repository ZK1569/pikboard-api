package repository

import (
	"errors"
	"strings"

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
	sess := self.db.DB.Session(&gorm.Session{})
	sess.Config.TranslateError = false

	user := model.User{
		Username: username,
		Email:    email,
		Password: password[:],
		Image:    "https://pikboard.s3.eu-west-1.amazonaws.com/profile_image/default_profile_image.jpg",
	}

	result := sess.Create(&user)
	if result.Error != nil {
		err := result.Error

		if strings.Contains(err.Error(), "duplicate key value") {
			if strings.Contains(err.Error(), "uni_users_username") {
				return nil, errs.UserAlreadyExists
			}
			if strings.Contains(err.Error(), "uni_users_email") {
				return nil, errs.EmailAlreadyExists
			}
			return nil, errs.AlreadyExists
		}

		return nil, err
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

func (self *User) GetUserByToken(token string) (*model.User, error) {
	var user model.User

	result := self.db.DB.
		Model(&model.User{}).
		Preload("Friends").
		Where("session = ?", token).
		First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.NotFound
		}
		return nil, result.Error
	}

	return &user, nil
}

func (self *User) GetUserByID(userID uint) (*model.User, error) {
	var user model.User

	result := self.db.DB.
		Model(&model.User{}).
		Preload("Friends").
		Where("id = ?", userID).
		First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.NotFound
		}
		return nil, result.Error
	}

	return &user, nil

}

func (self *User) SearchUsersByUsername(username string) ([]model.User, error) {
	users := []model.User{}

	result := self.db.DB.Where("username ILIKE ?", "%"+username+"%").Find(&users)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.NotFound
		}
		return nil, result.Error
	}

	return users, nil

}

func (self *User) UpdateUser(user *model.User) error {
	result := self.db.DB.Save(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (self *User) UpdatePassword(user *model.User, newPassword [64]byte) error {
	user.Password = newPassword[:]

	result := self.db.DB.Save(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
