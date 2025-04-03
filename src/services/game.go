package service

import (
	"github.com/google/uuid"
	errs "github.com/zk1569/pikboard-api/src/errors"
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

func (self *Game) CreateGame(owner *model.User, opponent *model.User, whitePlayerID uint, fem string) (*model.Game, error) {
	status, err := self.statusRepository.GetByStatus(model.StatusPending)
	if err != nil {
		return nil, err
	}
	game, err := self.gameRepository.CreateGame(owner.ID, opponent.ID, whitePlayerID, fem, status.ID)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (self *Game) GetUsersCurrentGame(user *model.User) ([]model.Game, error) {
	currentGames := []model.Game{}

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

func (self *Game) GetUsersRequestedGame(user *model.User) ([]model.Game, error) {
	requestedGames := []model.Game{}

	games, err := self.gameRepository.GetUsersGame(user.ID)
	if err != nil {
		return nil, err
	}

	for _, game := range games {
		if game.Status.Status == model.StatusPending {
			requestedGames = append(requestedGames, game)
		}
	}

	return requestedGames, nil
}

func (self *Game) GetUsersEndedGame(user *model.User) ([]model.Game, error) {
	endedGames := []model.Game{}

	games, err := self.gameRepository.GetUsersGame(user.ID)
	if err != nil {
		return nil, err
	}

	for _, game := range games {
		if game.Status.Status == model.StatusEnd {
			endedGames = append(endedGames, game)
		}
	}

	return endedGames, nil
}

func (self *Game) AcceptOrNotGame(gameID uint, user *model.User, answer bool) error {
	game, err := self.gameRepository.GetById(gameID)
	if err != nil {
		return err
	}

	if game.Opponent.ID != user.ID {
		return errs.Unauthorized
	}

	if answer {
		return self.acceptGameRequest(game)
	} else {
		return self.declineGameRequest(game)
	}
}

func (self *Game) acceptGameRequest(game *model.Game) error {

	status, err := self.statusRepository.GetByStatus(model.StatusPlaying)
	if err != nil {
		return err
	}

	game.Status = *status

	return self.gameRepository.Update(game)
}

func (self *Game) declineGameRequest(game *model.Game) error {
	return self.gameRepository.DeleteGame(game.ID)
}
