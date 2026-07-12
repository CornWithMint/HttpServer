package main

import (
	"net/http"
	"server/delivery"
	"server/repository"
	"server/usecase"
)

func main() {
	taskStorage := repository.NewStorage()
	userStorage := repository.NewUserStorage()

	taskUsecase := usecase.NewTaskUsecase(taskStorage)
	userUsecase := usecase.NewAuthUsecase(userStorage)

	mux := http.NewServeMux()

	delivery.Handlers(mux, taskUsecase, userUsecase)

	middlewars := delivery.RecoveryMiddleware(delivery.LoggingMiddleware(mux))

	server := http.Server{
		Addr:    ":8080",
		Handler: middlewars,
	}

	server.ListenAndServe()
}
