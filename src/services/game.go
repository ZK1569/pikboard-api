package service

import (
	"github.com/google/uuid"
	model "github.com/zk1569/pikboard-api/src/models"
	repository "github.com/zk1569/pikboard-api/src/repositories"
)

type Game struct {
	imageRepository  repository.ImageInterface
	gameRepository   repository.GameInterface
	statusRepository repository.StatusInterface
	iaRepository     repository.IAInterface
}

var singleGameInstance GameInterface

func GetGameInsance() GameInterface {
	if singleGameInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleGameInstance == nil {
			singleGameInstance = &Game{
				imageRepository:  repository.GetImageInstance(),
				iaRepository:     repository.GetChatGPTInstance(),
				gameRepository:   repository.GetGameInstance(),
				statusRepository: repository.GetStatusInstance(),
			}
		}
	}

	return singleGameInstance
}

func (self *Game) ImageToFem(img []byte) (string, error) {

	image_id := uuid.New().String()

	url, err := self.imageRepository.UploadForChat(image_id, img)
	if err != nil {
		return "", err
	}
	response, err := self.iaRepository.ImageToFem(url)
	if err != nil {
		return "", err
	}

	return response, nil
}

func (self *Game) CreateGame(owner *model.User, opponent *model.User, fem string) (*model.Game, error) {
	status, err := self.statusRepository.GetByStatus(model.StatusPending)
	if err != nil {
		return nil, err
	}
	game, err := self.gameRepository.CreateGame(owner.ID, opponent.ID, fem, status.ID)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (self *Game) GetUsersCurrentGame(user *model.User) ([]model.Game, error) {
	var currentGames []model.Game

	games, err := self.gameRepository.GetUsersGame(user.ID)
	if err != nil {
		return nil, err
	}

	for _, game := range games {
		if game.Status.Status == model.StatusPlaying {
			currentGames = append(currentGames, game)
		}
	}

	return currentGames, nil
}
