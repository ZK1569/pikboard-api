package controller

import (
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
		r.Post("/request", self.sendFriendRequest)
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

	_, err = self.friendService.SendFriendRequest(user, receiverUserID)
	if err != nil {
		jsonResponseError(w, err)
		return
	}

	jsonResponse(w, http.StatusCreated, nil)
}
