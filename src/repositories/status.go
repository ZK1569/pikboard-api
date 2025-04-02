package repository

import (
	"errors"

	errs "github.com/zk1569/pikboard-api/src/errors"
	model "github.com/zk1569/pikboard-api/src/models"
	util "github.com/zk1569/pikboard-api/src/utils"
	"gorm.io/gorm"
)

type Status struct {
	db *util.DatabasePostgres
}

var singleStatusInstance StatusInterface

func GetStatusInstance() StatusInterface {
	if singleStatusInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleStatusInstance == nil {
			singleStatusInstance = &Status{
				db: util.GetDatabasePostgresInstance(),
			}
		}
	}

	return singleStatusInstance
}

func (self *Status) CreateStatus(status string) (*model.Status, error) {
	var s = model.Status{
		Status: status,
	}

	r := self.db.DB.Create(&s)
	if r.Error != nil {
		if errors.Is(r.Error, gorm.ErrDuplicatedKey) {
			return nil, errs.AlreadyExists
		}
		return nil, r.Error
	}

	return &s, nil
}

func (self *Status) GetById(id uint) (*model.Status, error) {
	var s model.Status

	r := self.db.DB.First(&s, "id = ?", id)
	if r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			return nil, errs.NotFound
		}
		return nil, r.Error
	}

	return &s, nil
}

func (self *Status) GetByStatus(status string) (*model.Status, error) {
	var s model.Status

	r := self.db.DB.First(&s, "status = ?", status)
	if r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			return nil, errs.NotFound
		}
		return nil, r.Error
	}

	return &s, nil
}
