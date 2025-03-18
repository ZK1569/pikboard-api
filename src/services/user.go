package service

import (
	"crypto/sha512"

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
