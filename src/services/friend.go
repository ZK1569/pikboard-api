package service

import (
	"log"

	errs "github.com/zk1569/pikboard-api/src/errors"
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

	isFriend := self.IsFriend(senderUser, receiverUserID)
	if isFriend {
		return nil, errs.BadRequest
	}

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

func (self *Friend) AcceptOrNotFriendRequest(user *model.User, friend_id uint, answer bool) error {
	friend, err := self.userRepository.GetUserByID(friend_id)
	if err != nil {
		return err
	}

	if answer {
		return self.acceptFriendRequest(user, friend)
	} else {
		return self.declineFriendRequest(user, friend)
	}
}

func (self *Friend) acceptFriendRequest(user *model.User, friend *model.User) error {
	_, err := self.friendRepository.GetFriendRequest(user.ID, friend.ID)
	if err != nil {
		return err
	}

	user.Friends = append(user.Friends, friend)
	friend.Friends = append(friend.Friends, user)

	err = self.userRepository.UpdateUser(user)
	if err != nil {
		log.Printf("Error will updating user : %d", user.ID)
		return err
	}
	err = self.userRepository.UpdateUser(friend)
	if err != nil {
		log.Printf("Error will updating friend : %d", friend.ID)
		return err
	}

	err = self.friendRepository.DeleteFriendRequest(user.ID, friend.ID)
	if err != nil {
		return err
	}

	return nil
}

func (self *Friend) declineFriendRequest(user *model.User, friend *model.User) error {
	err := self.friendRepository.DeleteFriendRequest(user.ID, friend.ID)
	if err != nil {
		return err
	}

	return nil
}

func (self *Friend) IsFriend(user *model.User, friendID uint) bool {
	for _, i := range user.Friends {
		if i.ID == friendID {
			return true
		}
	}

	return false
}

func (self *Friend) GetPendingFriendRequest(user *model.User) ([]model.User, error) {
	friendRequests, err := self.friendRepository.GetPendingFriendRequest(user.ID)
	if err != nil {
		return nil, err
	}

	users := []model.User{}

	for _, request := range friendRequests {
		users = append(users, request.Requester)
	}

	return users, nil
}

func (self *Friend) GetSentFriendRequest(user *model.User) ([]model.User, error) {
	friendSenderd, err := self.friendRepository.GetSentFriendRequest(user.ID)
	if err != nil {
		return nil, err
	}

	users := []model.User{}
	for _, request := range friendSenderd {
		users = append(users, request.Receiver)
	}

	return users, nil
}
