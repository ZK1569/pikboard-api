package controller

import (
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	errs "github.com/zk1569/pikboard-api/src/errors"
	service "github.com/zk1569/pikboard-api/src/services"
)

type Game struct {
	path        string
	gameService service.GameInterface
}

var singleGameInstance *Game

func GetGameInsance() *Game {
	if singleGameInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleGameInstance == nil {
			singleGameInstance = &Game{
				path:        "/game",
				gameService: service.GetGameInsance(),
			}
		}
	}

	return singleGameInstance
}

func (self *Game) Mount(r chi.Router) {
	r.Route(self.path, func(r chi.Router) {
		r.Use(GetMiddlewareInstance().AuthTokenMiddleware)
		r.Post("/position", self.getPossitionFromImg)
	})
}

func (self *Game) getPossitionFromImg(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(100 << 20)

	if err != nil {
		jsonResponseError(w, err)
		return
	}

	file, handler, err := r.FormFile("img")
	if err != nil {
		jsonResponseError(w, err)
		return
	}
	defer file.Close()

	namesplit := strings.Split(handler.Filename, ".")
	file_extension := namesplit[len(namesplit)-1]

	if file_extension == "jpg" || file_extension == "png" {

		imageData, err := io.ReadAll(file)

		if err != nil {
			jsonResponseError(w, err)
			return
		}

		fem, err := self.gameService.ImageToFem(imageData)

		if err != nil {
			jsonResponseError(w, err)
			return
		}

		jsonResponse(w, http.StatusOK, fem)
		return
	}

	jsonResponseError(w, errs.BadRequest)

}
