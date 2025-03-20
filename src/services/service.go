package service

import (
	"sync"

	model "github.com/zk1569/pikboard-api/src/models"
)

var lock *sync.Mutex

func init() {
	lock = &sync.Mutex{}
}

type UserInterface interface {
	CreateUser(username, email, password string) (*model.User, error)
	GetUserSession(email, password string) (string, error)
	GetUserByID(uint) (*model.User, error)
	UpdateUser(*model.User) error
	UpdatePassword(user *model.User, oldPassword, newPassword string) error
	SearchUsersByUsername(string) ([]model.User, error)
	GetUserByToken(string) (*model.User, error)
}

type FriendInterface interface {
	SendFriendRequest(*model.User, uint) (*model.FriendRequest, error)
	AcceptOrNotFriendRequest(*model.User, uint, bool) error
	acceptFriendRequest(*model.User, *model.User) error
	declineFriendRequest(*model.User, *model.User) error
}
