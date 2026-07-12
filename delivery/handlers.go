package delivery

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"server/domain"
	"server/usecase"
	"strconv"
)

var UserExists = errors.New("User already exists")
var UserNotExists = errors.New("User not exists")
var BadRequest = errors.New("Bad request")
var InternalServerError = errors.New("Internal Server Error")
var ErrTaskNotFound = errors.New("task not found")
var ErrEmptyTitle = errors.New("Empty title")

func Handlers(mux *http.ServeMux, taskStorage *usecase.TaskUsecase, userStorage *usecase.AuthUsecase) {
	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		ans := &domain.Answer{
			Message: "pong",
		}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(ans)
		if err != nil {
			log.Println(err)
		}
	})

	mux.HandleFunc("POST /echo", func(w http.ResponseWriter, r *http.Request) {
		ans := &domain.Answer{}

		err := json.NewDecoder(r.Body).Decode(ans)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(ans)
		if err != nil {
			log.Println(err)
		}
	})

	mux.HandleFunc("GET /hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		ans := &domain.Answer{
			Hello: name,
		}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(ans)
		if err != nil {
			log.Println(err)
		}
	})

	mux.HandleFunc("GET /tasks*", func(w http.ResponseWriter, r *http.Request) {
		res, err := taskStorage.ListTasks()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Println(err)
		}
	})

	mux.HandleFunc("POST /tasks*", func(w http.ResponseWriter, r *http.Request) {
		task := &domain.Task{}
		err := json.NewDecoder(r.Body).Decode(task)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		res, err := taskStorage.CreateTask(task.Title)
		if err != nil {
			if errors.Is(err, ErrEmptyTitle) {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Println(err)
		}
	})

	mux.HandleFunc("GET /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		int_id, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		res, err := taskStorage.GetTaskByID(int_id)
		if err != nil {
			if errors.Is(err, ErrTaskNotFound) {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Println(err)
		}
	})

	mux.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
		user := &domain.User{}
		err := json.NewDecoder(r.Body).Decode(user)

		//Usecase возвращает User обязательно без пароля
		user, err = userStorage.Register(user.Username, user.Password)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			log.Println(err)
		}
	})

	mux.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
		user := &domain.User{}
		err := json.NewDecoder(r.Body).Decode(user)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		//Usecase возвращает token

		w.Header().Set("Content-Type", "application/json")
		// err = json.NewEncoder(w).Encode()
		// if err != nil {
		// 	log.Println(err)
		// }
	})
}
