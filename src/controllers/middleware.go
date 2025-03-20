package controller

import (
	"context"
	"log"
	"net/http"
	"strings"

	errs "github.com/zk1569/pikboard-api/src/errors"
	service "github.com/zk1569/pikboard-api/src/services"
)

type Middleware struct {
	userService service.UserInterface
}

var singleMiddlewareInstance *Middleware

func GetMiddlewareInstance() *Middleware {
	if singleMiddlewareInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleMiddlewareInstance == nil {
			singleMiddlewareInstance = &Middleware{
				userService: service.GetUserInstance(),
			}
		}
	}

	return singleMiddlewareInstance
}

func (self *Middleware) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Printf("authorization header is missing")
			jsonResponseError(w, errs.Unauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Printf("authorization header is malformed")
			jsonResponseError(w, errs.Unauthorized)
			return
		}

		token := parts[1]

		user, err := self.userService.GetUserByToken(token)
		if err != nil {
			jsonResponseError(w, errs.Unauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, userCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
