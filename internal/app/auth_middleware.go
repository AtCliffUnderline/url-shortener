package app

import (
	"context"
	"net/http"
)

const COOKIE_NAME string = "auth-token"

func (service *ShortenerService) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(COOKIE_NAME)
		if err != http.ErrNoCookie {
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			user, err := service.UserRepository.GetUserByToken(cookie.Value)
			// Если логин прошел успешно
			if err == nil {
				ctx := context.WithValue(r.Context(), "user", user)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}
		user := service.UserRepository.CreateUser()
		http.SetCookie(w, &http.Cookie{
			Name:  COOKIE_NAME,
			Value: user.Token,
		})
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
