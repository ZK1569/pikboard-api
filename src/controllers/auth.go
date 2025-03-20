package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	errs "github.com/zk1569/pikboard-api/src/errors"
	model "github.com/zk1569/pikboard-api/src/models"
	service "github.com/zk1569/pikboard-api/src/services"
)

type Authentification struct {
	path        string
	userService service.UserInterface
}

var singleAuthentificationInstance *Authentification

func GetAuthentificationInstance() *Authentification {
	if singleAuthentificationInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleAuthentificationInstance == nil {
			singleAuthentificationInstance = &Authentification{
				path:        "/",
				userService: service.GetUserInstance(),
			}
		}
	}

	return singleAuthentificationInstance
}

func (self *Authentification) Mount(r chi.Router) {
	r.Route(self.path, func(r chi.Router) {
		r.Post("/login", self.login)
		r.Post("/signup", self.signup)
	})
}

type LoginBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"gte=8,lte=100"`
}

func (self *Authentification) login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var bodyLogin LoginBody
	if err := json.NewDecoder(r.Body).Decode(&bodyLogin); err != nil {
		log.Printf("Error: body decode %v", err)
		jsonResponseError(w, errs.BadRequest)
		return
	}

	if err := Validate.Struct(bodyLogin); err != nil {
		log.Printf("Error: Validation error %v", err)
		jsonResponseError(w, errs.BadRequest)
		return
	}

	session_token, err := self.userService.GetUserSession(bodyLogin.Email, bodyLogin.Password)
	if err != nil {
		jsonResponseError(w, err)
		return
	}

	type tokenResponse struct {
		Token string `json:"token"`
	}
	jsonResponse(w, http.StatusOK, &tokenResponse{Token: session_token})
	return
}

type NewUserBody struct {
	Username string `json:"username" validate:"gte=5,lte=130"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"gte=8,lte=100"`
}

func (self *Authentification) signup(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var bodyUser NewUserBody
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

	user, err := self.userService.CreateUser(bodyUser.Username, bodyUser.Email, bodyUser.Password)
	if err != nil {
		jsonResponseError(w, err)
		return
	}

	type userResponse struct {
		User *model.User `json:"user"`
	}

	jsonResponse(w, http.StatusCreated, &userResponse{User: user})
	return
}
