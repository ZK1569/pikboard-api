package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	errs "github.com/zk1569/pikboard-api/src/errors"
	service "github.com/zk1569/pikboard-api/src/services"
)

type Game struct {
	path        string
	gameService service.GameInterface
	userService service.UserInterface
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
				userService: service.GetUserInstance(),
			}
		}
	}

	return singleGameInstance
}

func (self *Game) Mount(r chi.Router) {
	r.Route(self.path, func(r chi.Router) {
		r.Use(GetMiddlewareInstance().AuthTokenMiddleware)
		r.Post("/position", self.getPossitionFromImg)
		r.Post("/new", self.createNewGame)
	})
}

func (self *Game) getPossitionFromImg(w http.ResponseWriter, r *http.Request) {
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

type NewGameBody struct {
	Fem        string `json:"fem" validate:"gte=5"`
	OpponentID uint   `json:"opponent_id"`
}

func (self *Game) createNewGame(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	defer r.Body.Close()

	var bodyGame NewGameBody
	if err := json.NewDecoder(r.Body).Decode(&bodyGame); err != nil {
		log.Printf("Error: body decode %v", err)
		jsonResponseError(w, errs.BadRequest)
		return
	}

	if err := Validate.Struct(bodyGame); err != nil {
		log.Printf("Error: Validation error %v", err)
		jsonResponseError(w, errs.BadRequest)
		return
	}

	opponent, err := self.userService.GetUserByID(bodyGame.OpponentID)
	if err != nil {
		jsonResponseError(w, err)
		return
	}

	_, err = self.gameService.CreateGame(user, opponent, bodyGame.Fem)
	if err != nil {
		jsonResponseError(w, err)
		return
	}

	jsonResponse(w, http.StatusOK, nil)
}
