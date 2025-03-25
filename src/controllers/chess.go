package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	service "github.com/zk1569/pikboard-api/src/services"
)

type Chess struct {
	path         string
	chessService service.ChessInterface
}

var singleChessInstance *Chess

func GetChessInstance() *Chess {
	if singleChessInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleChessInstance == nil {
			singleChessInstance = &Chess{
				path:         "/chess",
				chessService: service.GetChessInstance(),
			}
		}
	}

	return singleChessInstance
}

func (self *Chess) Mount(r chi.Router) {
	r.Route(self.path, func(r chi.Router) {
		r.Get("/", self.femToImage)
	})
}

func (self *Chess) femToImage(w http.ResponseWriter, r *http.Request) {
	chessFEM := r.URL.Query().Get("q")
	chessPovStr := r.URL.Query().Get("pov")

	chessPov := true
	if chessPovStr == "black" {
		chessPov = false
	}

	chessImage, err := self.chessService.FemToImage(chessFEM, chessPov)

	if err != nil {
		jsonResponseError(w, err)
		return
	}

	imageResponse(w, http.StatusOK, chessImage)
}
