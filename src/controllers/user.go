package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

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
		r.Get("/friends", self.getUserFriend)
		r.Put("/password", self.updatePassword)
		r.Post("/image", self.updateProfileImage)
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

type UpdatePasswordBody struct {
	OldPassword string `json:"old_password" validate:"gte=8,lte=100"`
	NewPassword string `json:"new_password" validate:"gte=8,lte=100"`
}

func (self *User) updatePassword(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var bodyPassword UpdatePasswordBody
	if err := json.NewDecoder(r.Body).Decode(&bodyPassword); err != nil {
		log.Printf("Error: body decode %v", err)
		jsonResponseError(w, errs.BadRequest)
		return
	}

	if err := Validate.Struct(bodyPassword); err != nil {
		log.Printf("Error: Validation error %v", err)
		jsonResponseError(w, errs.BadRequest)
		return
	}

	user := getUserFromCtx(r)

	err := self.userService.UpdatePassword(user, bodyPassword.OldPassword, bodyPassword.NewPassword)
	if err != nil {
		jsonResponseError(w, err)
		return
	}

	jsonResponse(w, http.StatusOK, user)
}

func (self *User) getUserFriend(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	jsonResponse(w, http.StatusOK, user.Friends)
}

func (self *User) updateProfileImage(w http.ResponseWriter, r *http.Request) {

	user := getUserFromCtx(r)

	err := r.ParseMultipartForm(100 << 20)

	if err != nil {
		jsonResponseError(w, err)
		return
	}

	file, handler, err := r.FormFile("profile_image")
	if err != nil {
		jsonResponseError(w, err)
		return
	}
	defer file.Close()

	fmt.Printf("Fichier reçu : %s (%d bytes) \n", handler.Filename, handler.Size)

	namesplit := strings.Split(handler.Filename, ".")
	file_extension := namesplit[len(namesplit)-1]

	if file_extension == "jpg" || file_extension == "png" {
		imageData, err := io.ReadAll(file)
		if err != nil {
			jsonResponseError(w, err)
			return
		}

		err = self.userService.UpdateProfileImage(user, imageData, file_extension)
		if err != nil {
			jsonResponseError(w, err)
			return
		}

		jsonResponse(w, http.StatusOK, nil)
		return
	}

	jsonResponseError(w, errs.BadRequest)
}
