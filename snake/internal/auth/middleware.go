package auth

import (
	"errors"
	"log"
	"net/http"
	jwt_setup "snake_service/pkg/jwt-setup"
	"strings"
)

type appHandler func(http.ResponseWriter, *http.Request) error

func Middleware(h appHandler) http.HandlerFunc {
	log.Println("got into auth middleware")
	return func(w http.ResponseWriter, r *http.Request) {
		var appErr *AppError
		headerVal := r.Header.Get("Authorization")
		if headerVal == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(ErrWrongToken.Marshal())
			return
		}

		authHeaderArr := strings.Split(headerVal, " ")
		if len(authHeaderArr) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(ErrWrongToken.Marshal())
			return
		}
		tokenString := authHeaderArr[1]
		_, err := jwt_setup.ParseToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(ErrWrongToken.Marshal())
			return
		}
		err = h(w, r)
		if err != nil {
			if errors.As(err, &appErr) {
				if errors.Is(err, ErrNotFound) {
					w.WriteHeader(http.StatusNotFound)
					w.Write(ErrNotFound.Marshal())
					return
				}
				err := err.(*AppError)
				w.WriteHeader(http.StatusBadRequest)
				w.Write(err.Marshal())
				return
			}
			w.WriteHeader(http.StatusTeapot)
			w.Write(systemError(err.Error()).Marshal())
			return
		}
	}
}

func NoAuthMiddleware(h appHandler) http.HandlerFunc {
	log.Println("got into middleware")
	return func(w http.ResponseWriter, r *http.Request) {
		var appErr *AppError
		err := h(w, r)
		if err != nil {
			if errors.As(err, &appErr) {
				if errors.Is(err, ErrNotFound) {
					w.WriteHeader(http.StatusNotFound)
					w.Write(ErrNotFound.Marshal())
					return
				}
				err := err.(*AppError)
				w.WriteHeader(http.StatusBadRequest)
				w.Write(err.Marshal())
				return
			}
			w.WriteHeader(http.StatusTeapot)
			w.Write(systemError(err.Error()).Marshal())
		}
	}
}
