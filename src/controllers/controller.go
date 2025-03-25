package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	errs "github.com/zk1569/pikboard-api/src/errors"
	model "github.com/zk1569/pikboard-api/src/models"
)

var (
	Validate *validator.Validate
	lock     *sync.Mutex
	userCtx  = "user"
)

type MangaInterface interface {
	Mount(r chi.Router)
}

type ChapterInterface interface {
	Mount(r chi.Router)
}

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
	lock = &sync.Mutex{}
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578 // 1 Mb
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func jsonResponseError(w http.ResponseWriter, err error) error {

	type envelope struct {
		Error string `json:"error"`
	}

	switch err {
	case errs.NotFound:
		return writeJSON(w, http.StatusNotFound, &envelope{Error: "Not found"})
	case errs.Validation, errs.BadRequest:
		return writeJSON(w, http.StatusBadRequest, &envelope{Error: "Bad Request"})
	case errs.AlreadyExists:
		return writeJSON(w, http.StatusConflict, &envelope{Error: "Already exists"})
	case errs.Unauthorized:
		return writeJSON(w, http.StatusUnauthorized, &envelope{Error: "Unauthorized"})
	default:
		log.Printf("❌ Error: %v \n", err)
		return writeJSON(w, http.StatusInternalServerError, &envelope{Error: "Internal error"})
	}

}

func jsonResponse(w http.ResponseWriter, status int, data any) error {
	type envelope struct {
		Data any `json:"data"`
	}

	return writeJSON(w, status, &envelope{Data: data})
}

func getUserFromCtx(r *http.Request) *model.User {
	user, _ := r.Context().Value(userCtx).(*model.User)
	return user
}

func SetupCORS(r chi.Router) {
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

}
