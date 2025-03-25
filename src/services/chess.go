package service

import (
	"errors"

	client "github.com/zk1569/pikboard-api/src/clients"
	errs "github.com/zk1569/pikboard-api/src/errors"
)

type Chess struct {
	chessImageClient client.ChessImage
}

var singleChessInstance ChessInterface

func GetChessInstance() ChessInterface {
	if singleChessInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleChessInstance == nil {
			singleChessInstance = &Chess{
				chessImageClient: client.GetChessVisionInstance(),
			}
		}
	}

	return singleChessInstance
}

func (self *Chess) FemToImage(fem string, isWhitePOV bool) ([]byte, error) {
	chessImage, err := self.chessImageClient.FemToImage(fem, isWhitePOV)

	if errors.Is(err, errs.ClientResponseNoOK) {
		return chessImage, errs.BadRequest
	}

	return chessImage, err
}
