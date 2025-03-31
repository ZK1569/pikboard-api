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

func (self *Friend) DeleteFriendRequest(user, friend uint) error {
	friendRequest := model.FriendRequest{}

	result := self.db.DB.Where(
		"(requester_id = ? AND receiver_id = ?) OR (requester_id = ? AND receiver_id = ?)",
		user,
		friend,
		friend,
		user,
	).Delete(&friendRequest)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (self *Friend) GetFriendRequest(user, friend uint) (*model.FriendRequest, error) {
	friendRequest := model.FriendRequest{}

	result := self.db.DB.Where(
		"requester_id = ? AND receiver_id = ?",
		friend,
		user,
	).First(&friendRequest)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.NotFound
		}
	}

	return &friendRequest, nil
}

func (self *Friend) GetPendingFriendRequest(userID uint) ([]model.FriendRequest, error) {
	friendRequests := []model.FriendRequest{}

	result := self.db.DB.Model(&model.FriendRequest{}).
		Preload("Receiver").
		Preload("Requester").
		Where(
			"receiver_id = ?",
			userID,
		).
		Find(&friendRequests)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.NotFound
		}
		return nil, result.Error
	}

	return friendRequests, nil
}

func (self *Friend) GetSentFriendRequest(userID uint) ([]model.FriendRequest, error) {
	friendRequests := []model.FriendRequest{}

	result := self.db.DB.Model(&model.FriendRequest{}).
		Preload("Receiver").
		Preload("Requester").
		Where(
			"requester_id = ?",
			userID,
		).
		Find(&friendRequests)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.NotFound
		}
		return nil, result.Error
	}

	return friendRequests, nil
}
