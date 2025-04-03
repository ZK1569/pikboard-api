package repository

import (
	"sync"

	model "github.com/zk1569/pikboard-api/src/models"
)

var lock *sync.Mutex

func init() {
	lock = &sync.Mutex{}
}

type UserInterface interface {
	CreateUser(username, email string, password [64]byte) (*model.User, error)
	GetUserByEmailAndPassword(string, [64]byte) (*model.User, error)
	GetUserByID(uint) (*model.User, error)
	SearchUsersByUsername(string) ([]model.User, error)
	UpdateUser(*model.User) error
	UpdateUserSession(*model.User, string) error
	UpdatePassword(user *model.User, newPassword [64]byte) error
	GetUserByToken(string) (*model.User, error)
}

type FriendInterface interface {
	CreateFriendRequest(*model.User, *model.User) (*model.FriendRequest, error)
	DeleteFriendRequest(uint, uint) error
	GetFriendRequest(uint, uint) (*model.FriendRequest, error)

	GetPendingFriendRequest(uint) ([]model.FriendRequest, error)
	GetSentFriendRequest(uint) ([]model.FriendRequest, error)
}

type ImageInterface interface {
	UploadImage(string, []byte, string) (string, error)
	UploadForChat(string, []byte) (string, error)
}

type IAInterface interface {
	ImageToFem(string) (string, error)
}

type GameInterface interface {
	CreateGame(uint, uint, uint, string, uint) (*model.Game, error)
	GetUsersGame(uint) ([]model.Game, error)
	GetById(uint) (*model.Game, error)
	DeleteGame(uint) error
	Update(*model.Game) error
}

type StatusInterface interface {
	CreateStatus(string) (*model.Status, error)
	GetById(uint) (*model.Status, error)
	GetByStatus(string) (*model.Status, error)
}
