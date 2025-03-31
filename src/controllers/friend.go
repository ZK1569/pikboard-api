package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	errs "github.com/zk1569/pikboard-api/src/errors"
	service "github.com/zk1569/pikboard-api/src/services"
)

type Friend struct {
	path          string
	userService   service.UserInterface
	friendService service.FriendInterface
}

var singleFriendInstance *Friend

func GetFriendInstance() *Friend {
	if singleFriendInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleFriendInstance == nil {
			singleFriendInstance = &Friend{
				path:          "/friend",
				userService:   service.GetUserInstance(),
				friendService: service.GetFriendInstance(),
			}
		}
	}

	return singleFriendInstance
}

func (self *Friend) Mount(r chi.Router) {
	r.Route(self.path, func(r chi.Router) {
		r.Use(GetMiddlewareInstance().AuthTokenMiddleware)
		r.Get("/request", self.getPendingFriendRequest)
		r.Get("/sent", self.getSentFriendRequest)
		r.Post("/request", self.sendFriendRequest)
		r.Post("/accept", self.acceptFriendRequest)
	})
}

func (self *Friend) sendFriendRequest(w http.ResponseWriter, r *http.Request) {

	user := getUserFromCtx(r)

	receiverUserIDstr := r.URL.Query().Get("id")
	if receiverUserIDstr == "" {
		jsonResponseError(w, errs.BadRequest)
		return
	}
	userID64, err := strconv.ParseInt(receiverUserIDstr, 10, 32)
	if err != nil {
		jsonResponseError(w, errs.BadRequest)
		return
	}
	receiverUserID := uint(userID64)

	if user.ID == receiverUserID {
		jsonResponseError(w, errs.BadRequest)
		return
	}

	_, err = self.friendService.SendFriendRequest(user, receiverUserID)
	if err != nil {
		jsonResponseError(w, err)
		return
	}

	jsonResponse(w, http.StatusCreated, nil)
}

type AcceptOrNotBody struct {
	Answer bool `json:"answer" validate:"required"`
}

func (self *Friend) acceptFriendRequest(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	user := getUserFromCtx(r)

	friendRequestIDstr := r.URL.Query().Get("friend_id")
	if friendRequestIDstr == "" {
		jsonResponseError(w, errs.BadRequest)
		return
	}
	userID64, err := strconv.ParseInt(friendRequestIDstr, 10, 32)
	if err != nil {
		jsonResponseError(w, errs.BadRequest)
		return
	}
	friendRequestID := uint(userID64)

	var bodyAnswer AcceptOrNotBody
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

	err = self.friendService.AcceptOrNotFriendRequest(user, friendRequestID, bodyAnswer.Answer)
	if err != nil {
		jsonResponseError(w, err)
		return
	}

	jsonResponse(w, http.StatusOK, nil)
}

func (self *Friend) getPendingFriendRequest(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	users, err := self.friendService.GetPendingFriendRequest(user)
	if err != nil {
		jsonResponseError(w, err)
		return
	}

	jsonResponse(w, http.StatusOK, users)
}

func (self *Friend) getSentFriendRequest(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	users, err := self.friendService.GetSentFriendRequest(user)
	if err != nil {
		jsonResponseError(w, err)
		return
	}

	jsonResponse(w, http.StatusOK, users)
}
