package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	util "github.com/zk1569/pikboard-api/src/utils"
)

type Health struct {
	path string
}

var singleHealthInstance *Health

func GetHealthInstance() *Health {
	if singleHealthInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		singleHealthInstance = &Health{
			path: "/health",
		}
	}

	return singleHealthInstance
}

func (self *Health) Mount(r chi.Router) {
	r.Route(self.path, func(r chi.Router) {
		r.Get("/", self.heathCheck)
	})

}

func (self *Health) heathCheck(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "OK",
		"env":     util.EnvVariable.Env,
		"version": util.EnvVariable.Version,
	}

	jsonResponse(w, http.StatusOK, data)
}
