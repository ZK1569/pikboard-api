package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	controller "github.com/zk1569/pikboard-api/src/controllers"
)

func main() {
	adr := "0.0.0.0:8080"
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)

	controller.SetupCORS(mux)

	mux.Use(middleware.Timeout(10 * time.Second))

	mux.Route("/v1", func(r chi.Router) {
		controller.GetHealthInstance().Mount(r)
		controller.GetAuthentificationInstance().Mount(r)
		controller.GetUserInstance().Mount(r)
		controller.GetFriendInstance().Mount(r)
		controller.GetChessInstance().Mount(r)
		controller.GetGameInsance().Mount(r)
	})

	srv := http.Server{
		Addr:         adr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server is running at %s\n", adr)
	srv.ListenAndServe()

}
