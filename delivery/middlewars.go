package delivery

import (
	"context"
	"log"
	"net/http"
	"server/domain"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("Path %s, Method %s, Time %s", r.URL.Path, r.Method, time.Since(start))

	})
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		parts := strings.Split(authHeader, " ")
		if !strings.EqualFold(parts[0], "bearer") {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		if parts[1] == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		claims := &domain.CustomClaims{}

		_, err := jwt.ParseWithClaims(parts[1], claims, func(t *jwt.Token) (any, error) {
			return nil, nil
		})
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", claims.Userid)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
