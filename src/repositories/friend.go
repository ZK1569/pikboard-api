package repository

import (
	"errors"

	errs "github.com/zk1569/pikboard-api/src/errors"
	model "github.com/zk1569/pikboard-api/src/models"
	util "github.com/zk1569/pikboard-api/src/utils"
	"gorm.io/gorm"
)

type Friend struct {
	db *util.DatabasePostgres
}

var singleFriendInstance FriendInterface

func GetFriendInstance() FriendInterface {
	if singleFriendInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleFriendInstance == nil {
			singleFriendInstance = &Friend{
				db: util.GetDatabasePostgresInstance(),
			}
		}
	}

	return singleFriendInstance
}

func (self *Friend) CreateFriendRequest(senderUser *model.User, receiverUser *model.User) (*model.FriendRequest, error) {

	friendRequest := model.FriendRequest{
		RequesterID: senderUser.ID,
		ReceiverID:  receiverUser.ID,
	}

	result := self.db.DB.Create(&friendRequest)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, errs.AlreadyExists
		}
		return nil, result.Error
	}

	return &friendRequest, nil

}
