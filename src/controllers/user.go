package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	service "github.com/zk1569/pikboard-api/src/services"
)

type User struct {
	path        string
	userService service.UserInterface
}

var singleUserInstance *User

func GetUserInstance() *User {
	if singleUserInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleUserInstance == nil {
			singleUserInstance = &User{
				path:        "/user",
				userService: service.GetUserInstance(),
			}
		}
	}

	return singleUserInstance
}

func (self *User) Mount(r chi.Router) {
	r.Route(self.path, func(r chi.Router) {
		r.Use(GetMiddlewareInstance().AuthTokenMiddleware)
		r.Get("/self", self.selfInfo)
	})
}

func (self *User) selfInfo(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	jsonResponse(w, http.StatusOK, user)
}
