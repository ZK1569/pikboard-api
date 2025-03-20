package service

import (
	"bytes"
	"crypto/sha512"

	"github.com/google/uuid"
	errs "github.com/zk1569/pikboard-api/src/errors"
	model "github.com/zk1569/pikboard-api/src/models"
	repository "github.com/zk1569/pikboard-api/src/repositories"
)

type User struct {
	userRepository repository.UserInterface
}

var singleUserInstance UserInterface

func GetUserInstance() UserInterface {
	if singleUserInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleUserInstance == nil {
			singleUserInstance = &User{
				userRepository: repository.GetUserInstance(),
			}
		}
	}

	return singleUserInstance
}

func (self *User) CreateUser(username, email, password string) (*model.User, error) {

	hash := sha512.Sum512([]byte(password))

	return self.userRepository.CreateUser(username, email, hash)
}

func (self *User) GetUserSession(email, password string) (string, error) {
	hash := sha512.Sum512([]byte(password))

	user, err := self.userRepository.GetUserByEmailAndPassword(email, hash)
	if err != nil {
		return "", errs.Unauthorized
	}

	sessionToken := uuid.New().String()

	err = self.userRepository.UpdateUserSession(user, sessionToken)
	if err != nil {
		return "", err
	}

	return sessionToken, nil
}

func (self *User) GetUserByToken(token string) (*model.User, error) {
	return self.userRepository.GetUserByToken(token)
}

func (self *User) GetUserByID(userID uint) (*model.User, error) {
	return self.userRepository.GetUserByID(userID)
}

func (self *User) SearchUsersByUsername(username string) ([]model.User, error) {
	return self.userRepository.SearchUsersByUsername(username)
}
func (self *User) UpdateUser(user *model.User) error {
	return self.userRepository.UpdateUser(user)
}

func (self *User) UpdatePassword(user *model.User, oldPassword, newPassword string) error {
	oldPasswordHash := sha512.Sum512([]byte(oldPassword))

	if !bytes.Equal(user.Password, oldPasswordHash[:]) {
		return errs.Unauthorized
	}

	newPasswordHash := sha512.Sum512([]byte(newPassword))

	err := self.userRepository.UpdatePassword(user, newPasswordHash)
	if err != nil {
		return err
	}

	return nil
}
