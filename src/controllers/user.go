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
		r.Get("/search", self.searchUser)
		r.Patch("/", self.updateUser)
	})
}

func (self *User) selfInfo(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	jsonResponse(w, http.StatusOK, user)
}

func (self *User) getUserByID(w http.ResponseWriter, r *http.Request) {
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

func (self *User) searchUser(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		jsonResponseError(w, errs.BadRequest)
		return
	}

	users, err := self.userService.SearchUsersByUsername(username)
	if err != nil {
		jsonResponseError(w, err)
		return
	}

	jsonResponse(w, http.StatusOK, users)
}

type UpdateUserBody struct {
	Email string `json:"email" validate:"omitempty,email"`
	Phone string `json:"phone"`
}

func (self *User) updateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var bodyUser UpdateUserBody
	if err := json.NewDecoder(r.Body).Decode(&bodyUser); err != nil {
		log.Printf("Error: body decode %v", err)
		jsonResponseError(w, errs.BadRequest)
		return
	}

	if err := Validate.Struct(bodyUser); err != nil {
		log.Printf("Error: Validation error %v", err)
		jsonResponseError(w, errs.BadRequest)
		return
	}

	user := getUserFromCtx(r)

	if bodyUser.Email != "" {
		user.Email = bodyUser.Email
	}
	if bodyUser.Phone != "" {
		user.Phone = &bodyUser.Phone
	}

	err := self.userService.UpdateUser(user)
	if err != nil {
		jsonResponseError(w, err)
		return
	}

	jsonResponse(w, http.StatusOK, user)
}
