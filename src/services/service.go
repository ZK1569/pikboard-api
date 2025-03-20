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
}
