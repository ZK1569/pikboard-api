package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	errs "github.com/zk1569/pikboard-api/src/errors"
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
		r.Get("/", self.getUserByID)
		r.Get("/self", self.selfInfo)
	})
}

func (self *User) selfInfo(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	jsonResponse(w, http.StatusOK, user)
}

func (self *User) getUserByID(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Query())
	userIDstr := r.URL.Query().Get("id")
	if userIDstr == "" {
		jsonResponseError(w, errs.BadRequest)
		return
	}
	userID64, err := strconv.ParseInt(userIDstr, 10, 32)
	if err != nil {
		jsonResponseError(w, errs.BadRequest)
		return
	}
	userID := uint(userID64)

	user, err := self.userService.GetUserByID(userID)
	if err != nil {
		jsonResponseError(w, err)
		return
	}

	jsonResponse(w, http.StatusOK, user)
}
