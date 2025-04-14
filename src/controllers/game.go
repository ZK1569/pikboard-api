package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5"
	errs "github.com/zk1569/pikboard-api/src/errors"
	service "github.com/zk1569/pikboard-api/src/services"
)

type Game struct {
	path        string
	gameService service.GameInterface
	userService service.UserInterface

	clients ClientList
	sync.RWMutex
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

				clients: make(ClientList),
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
		r.Post("/end", self.endGame)
		r.Post("/accept", self.acceptOrNotGame)
		r.Post("/position", self.getPossitionFromImg)
		r.Post("/new", self.createNewGame)
		r.HandleFunc("/chess", self.serverWS)
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

	game, err := self.gameService.CreateGame(user, opponent, bodyGame.WhitePlayerID, bodyGame.Fen)
	if err != nil {
		jsonResponseError(w, err)
		return
	}

	jsonResponse(w, http.StatusCreated, game)
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

type EndGameBody struct {
	GameID   int `json:"game_id" validate:"required"`
	WinnerID int `json:"winner_id" validate:"required"`
}

func (self *Game) endGame(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	user := getUserFromCtx(r)

	var gameBody EndGameBody
	if err := json.NewDecoder(r.Body).Decode(&gameBody); err != nil {
		log.Printf("Error: body decode %v", err)
		jsonResponseError(w, errs.BadRequest)
		return
	}

	if err := Validate.Struct(gameBody); err != nil {
		log.Printf("Error: Validation error %v", err)
		jsonResponseError(w, errs.BadRequest)
		return
	}

	isOwner, err := self.gameService.IsUserOwner(user, uint(gameBody.GameID))
	if err != nil {
		jsonResponseError(w, err)
		return
	}
	if !isOwner {
		jsonResponseError(w, errs.Unauthorized)
		return
	}

	err = self.gameService.EndGame(uint(gameBody.GameID), uint(gameBody.WinnerID))
	if err != nil {
		jsonResponseError(w, err)
		return
	}

	jsonResponse(w, http.StatusOK, nil)
}

func (self *Game) serverWS(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

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

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(conn, self, gameID)

	self.addClient(client)

	go client.readMessage()
	go client.writeMessage()

}

func (self *Game) addClient(client *Client) {
	self.Lock()
	defer self.Unlock()

	self.clients[client] = true
}

func (self *Game) removeClient(client *Client) {
	self.Lock()
	defer self.Unlock()

	if _, ok := self.clients[client]; ok {
		client.connection.Close()
		delete(self.clients, client)
	}
}

func (self *Game) playAMove(gameID uint, newPosition string) error {

	game, err := self.gameService.GetByID(gameID)
	if err != nil {
		log.Printf("Error a la reception du message jpc a la reception: %v", err)
		return err
	}

	self.gameService.MakeAMove(game, newPosition)

	self.sendBackToSameGame(gameID)

	return nil
}

func (self *Game) sendBackToSameGame(gameID uint) error {
	game, err := self.gameService.GetByID(gameID)
	if err != nil {
		log.Printf("Error a la reception du message jpc: %v", err)
		return err
	}

	for wsclient := range self.clients {
		if wsclient.gameID == gameID {
			wsclient.egress <- []byte(game.Board)
		}
	}
	return nil
}
