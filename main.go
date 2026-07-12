package main

import (
	"net/http"
	"server/delivery"
	"server/repository"
	"server/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {
	taskStorage := repository.NewStorage()
	userStorage := repository.NewUserStorage()

	taskUsecase := usecase.NewTaskUsecase(taskStorage)
	userUsecase := usecase.NewAuthUsecase(userStorage)

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	delivery.Handlers(r, taskUsecase, userUsecase)

	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	server.ListenAndServe()
}
