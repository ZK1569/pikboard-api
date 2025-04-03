package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
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
		r.Get("/current", self.getCurrentGames)
		r.Get("/request", self.getRequestedGames)
		r.Get("/end", self.getEndedGames)
		r.Post("/accept", self.acceptOrNotGame)
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
	Fen           string `json:"fen" validate:"gte=5"`
	OpponentID    uint   `json:"opponent_id"`
	WhitePlayerID uint   `json:"white_player_id"`
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

	if user.ID == bodyGame.OpponentID {
		jsonResponseError(w, errs.BadRequest)
		return
	}

	opponent, err := self.userService.GetUserByID(bodyGame.OpponentID)
	if err != nil {
		jsonResponseError(w, err)
		return
	}

	_, err = self.gameService.CreateGame(user, opponent, bodyGame.WhitePlayerID, bodyGame.Fen)
	if err != nil {
		jsonResponseError(w, err)
		return
	}

	jsonResponse(w, http.StatusCreated, nil)
}

func (self *Game) getCurrentGames(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	games, err := self.gameService.GetUsersCurrentGame(user)
	if err != nil {
		jsonResponseError(w, err)
		return
	}

	jsonResponse(w, http.StatusOK, games)
}

func (self *Game) getRequestedGames(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	games, err := self.gameService.GetUsersRequestedGame(user)
	if err != nil {
		jsonResponseError(w, err)
		return
	}

	jsonResponse(w, http.StatusOK, games)
}

func (self *Game) getEndedGames(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	games, err := self.gameService.GetUsersEndedGame(user)
	if err != nil {
		jsonResponseError(w, err)
		return
	}

	jsonResponse(w, http.StatusOK, games)

}

type AcceptOrNotGameBody struct {
	Answer bool `json:"answer"`
}

func (self *Game) acceptOrNotGame(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	user := getUserFromCtx(r)

	gameIDstr := r.URL.Query().Get("g")
	if gameIDstr == "" {
		jsonResponseError(w, errs.BadRequest)
		return
	}
	gameID64, err := strconv.ParseInt(gameIDstr, 10, 32)
	if err != nil {
		jsonResponseError(w, errs.BadRequest)
		return
	}
	gameID := uint(gameID64)

	var bodyAnswer AcceptOrNotGameBody
	if err := json.NewDecoder(r.Body).Decode(&bodyAnswer); err != nil {
		log.Printf("Error: body decode %v", err)
		jsonResponseError(w, errs.BadRequest)
		return
	}

	if err := Validate.Struct(bodyAnswer); err != nil {
		log.Printf("Error: Validation error %v", err)
		jsonResponseError(w, errs.BadRequest)
		return
	}

	err = self.gameService.AcceptOrNotGame(gameID, user, bodyAnswer.Answer || false)
	if err != nil {
		jsonResponseError(w, err)
		return
	}

	jsonResponse(w, http.StatusOK, nil)
}
