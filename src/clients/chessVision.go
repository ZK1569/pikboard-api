package client

import (
	"fmt"
	"io"
	"log"
	"net/http"

	errs "github.com/zk1569/pikboard-api/src/errors"
)

type ChessVision struct {
	url string
}

var singleChessVisionInstance ChessImage

func GetChessVisionInstance() ChessImage {
	if singleChessVisionInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleChessVisionInstance == nil {
			singleChessVisionInstance = &ChessVision{
				url: "https://fen2image.chessvision.ai/",
			}
		}
	}

	return singleChessVisionInstance
}

func (self *ChessVision) FemToImage(fem string, isWhitePov bool) ([]byte, error) {
	pov := ""
	if isWhitePov == true {
		pov = "white"
	} else {
		pov = "black"
	}

	url := fmt.Sprintf("%s%s?pov=%s", self.url, fem, pov)
	log.Println(url)

	response, err := http.Get(url)
	if err != nil {
		return nil, errs.ClientResponseNoOK
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errs.ClientResponseNoOK
	}

	if ct := response.Header.Get("Content-Type"); ct != "image/png" {
		return nil, errs.ClientNotRightType
	}

	imageData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Error when reading response body: %v", err)
	}

	return imageData, nil
}
