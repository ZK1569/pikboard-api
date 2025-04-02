package repository

import (
	"errors"

	errs "github.com/zk1569/pikboard-api/src/errors"
	model "github.com/zk1569/pikboard-api/src/models"
	util "github.com/zk1569/pikboard-api/src/utils"
	"gorm.io/gorm"
)

type Game struct {
	db *util.DatabasePostgres
}

var singleGameInstance GameInterface

func GetGameInstance() GameInterface {
	if singleGameInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleGameInstance == nil {
			singleGameInstance = &Game{
				db: util.GetDatabasePostgresInstance(),
			}
		}
	}

	return singleGameInstance
}

func (self *Game) CreateGame(ownerID uint, opponentID uint, fem string, statusID uint) (*model.Game, error) {
	game := model.Game{
		UserID:     ownerID,
		OpponentID: opponentID,
		Board:      fem,
		StatusID:   statusID,
	}

	result := self.db.DB.Create(&game)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, errs.AlreadyExists
		}
		return nil, result.Error
	}

	return &game, nil
}

func (self *Game) GetUsersGame(userID uint) ([]model.Game, error) {
	var games []model.Game

	r := self.db.DB.Where("user_id = ? or opponent_id = ?", userID, userID).Preload("User").Preload("Opponent").Preload("Status").Find(&games)
	if r.Error != nil {
		return nil, r.Error
	}

	return games, nil
}
