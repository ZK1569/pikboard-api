package service

import (
	"github.com/google/uuid"
	repository "github.com/zk1569/pikboard-api/src/repositories"
)

type Game struct {
	imageRepository repository.ImageInterface
	iaRepository    repository.IAInterface
}

var singleGameInstance GameInterface

func GetGameInsance() GameInterface {
	if singleGameInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleGameInstance == nil {
			singleGameInstance = &Game{
				imageRepository: repository.GetImageInstance(),
				iaRepository:    repository.GetChatGPTInstance(),
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
