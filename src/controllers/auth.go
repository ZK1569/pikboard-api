package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Authentification struct {
	path string
}

var singleAuthentificationInstance *Authentification

func GetAuthentificationInstance() *Authentification {
	if singleAuthentificationInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleAuthentificationInstance == nil {
			singleAuthentificationInstance = &Authentification{
				path: "/",
			}
		}
	}

	return singleAuthentificationInstance
}

func (self *Authentification) Mount(r chi.Router) {
	r.Route(self.path, func(r chi.Router) {
		r.Post("/login", self.login)
		r.Get("/signup", self.signup)
	})
}

func (self *Authentification) login(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusNotImplemented, nil)
}

func (self *Authentification) signup(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusNotImplemented, nil)
}
