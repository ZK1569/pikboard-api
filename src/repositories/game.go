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

func (self *Game) CreateGame(ownerID uint, opponentID uint, whitePlayerID uint, fem string, statusID uint) (*model.Game, error) {
	game := model.Game{
		UserID:        ownerID,
		OpponentID:    opponentID,
		Board:         fem,
		StatusID:      statusID,
		WhitePlayerID: &whitePlayerID,
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

func (self *Game) GetById(gameID uint) (*model.Game, error) {
	var game model.Game

	result := self.db.DB.Model(&model.Game{}).Preload("User").Preload("Opponent").Preload("Status").Where("id = ?", gameID).First(&game)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.NotFound
		}
		return nil, result.Error
	}

	return &game, nil
}

func (self *Game) DeleteGame(gameID uint) error {
	game := model.Game{}

	result := self.db.DB.Where("id = ?", gameID).Delete(&game)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (self *Game) Update(game *model.Game) error {
	result := self.db.DB.Save(&game)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
