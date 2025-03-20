package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	service "github.com/zk1569/pikboard-api/src/services"
)

type Friend struct {
	path        string
	userService service.UserInterface
}

var singleFriendInstance *Friend

func GetFriendInstance() *Friend {
	if singleFriendInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleFriendInstance == nil {
			singleFriendInstance = &Friend{
				path:        "/friend",
				userService: service.GetUserInstance(),
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
	jsonResponse(w, http.StatusNotImplemented, nil)
}
