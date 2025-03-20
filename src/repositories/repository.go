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
	UpdateUserSession(*model.User, string) error
	GetUserByToken(string) (*model.User, error)
}
