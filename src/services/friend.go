package service

import (
	model "github.com/zk1569/pikboard-api/src/models"
	repository "github.com/zk1569/pikboard-api/src/repositories"
)

type Friend struct {
	userRepository   repository.UserInterface
	friendRepository repository.FriendInterface
}

var singleFriendInstance FriendInterface

func GetFriendInstance() FriendInterface {
	if singleFriendInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleFriendInstance == nil {
			singleFriendInstance = &Friend{
				userRepository:   repository.GetUserInstance(),
				friendRepository: repository.GetFriendInstance(),
			}
		}
	}

	return singleFriendInstance
}

func (self *Friend) SendFriendRequest(senderUser *model.User, receiverUserID uint) (*model.FriendRequest, error) {

	receiverUser, err := self.userRepository.GetUserByID(receiverUserID)
	if err != nil {
		return nil, err
	}

	newFriendRequest, err := self.friendRepository.CreateFriendRequest(senderUser, receiverUser)
	if err != nil {
		return nil, err
	}

	return newFriendRequest, nil
}
