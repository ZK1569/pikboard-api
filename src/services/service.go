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
	UpdateProfileImage(*model.User, []byte, string) error
}

type FriendInterface interface {
	SendFriendRequest(*model.User, uint) (*model.FriendRequest, error)
	AcceptOrNotFriendRequest(*model.User, uint, bool) error
	acceptFriendRequest(*model.User, *model.User) error
	declineFriendRequest(*model.User, *model.User) error
	IsFriend(*model.User, uint) bool

	GetPendingFriendRequest(*model.User) ([]model.User, error)
	GetSentFriendRequest(*model.User) ([]model.User, error)
}

type ChessInterface interface {
	FemToImage(string, bool) ([]byte, error)
}

type GameInterface interface {
	ImageToFem([]byte) (string, error)
	CreateGame(*model.User, *model.User, uint, string) (*model.Game, error)
	GetUsersCurrentGame(*model.User) ([]model.Game, error)
	GetUsersRequestedGame(*model.User) ([]model.Game, error)
	GetUsersEndedGame(*model.User) ([]model.Game, error)
	AcceptOrNotGame(uint, *model.User, bool) error
	EndGame(uint, uint) error
	IsUserOwner(*model.User, uint) (bool, error)
	GetByID(uint) (*model.Game, error)
	MakeAMove(*model.Game, string) error
}
