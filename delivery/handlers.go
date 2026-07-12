package delivery

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"server/domain"
	"server/usecase"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var UserExists = errors.New("User already exists")
var UserNotExists = errors.New("User not exists")
var BadRequest = errors.New("Bad request")
var InternalServerError = errors.New("Internal Server Error")
var ErrTaskNotFound = errors.New("task not found")
var ErrEmptyTitle = errors.New("Empty title")
var ErrInvalidCredentials = errors.New("Err Invalid Credentials")

func Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ans := &domain.Answer{
			Message: "pong",
		}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(ans)
		if err != nil {
			log.Println(err)
		}
	}
}

func Echo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

func Hello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		ans := &domain.Answer{
			Hello: name,
		}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(ans)
		if err != nil {
			log.Println(err)
		}
	}

}

func Register(userStorage *usecase.AuthUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := &domain.User{}
		err := json.NewDecoder(r.Body).Decode(user)

		user, err = userStorage.Register(user.Username, user.Password)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		user.Password = ""
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			log.Println(err)
		}
	}

}

func Login(userStorage *usecase.AuthUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := &domain.User{}
		err := json.NewDecoder(r.Body).Decode(user)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		token, err := userStorage.Login(user.Username, user.Password)
		if err != nil {
			if errors.Is(err, BadRequest) {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			if errors.Is(err, ErrInvalidCredentials) {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(token)
		if err != nil {
			log.Println(err)
		}
	}

}

func GetTasks(taskStorage *usecase.TaskUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userid := r.Context().Value("userID")
		res, err := taskStorage.ListTasks(userid.(uuid.UUID))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Println(err)
		}
	}

}

func PostTasks(taskStorage *usecase.TaskUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		task := &domain.Task{}
		err := json.NewDecoder(r.Body).Decode(task)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		userID := r.Context().Value("userID")

		res, err := taskStorage.CreateTask(userID.(uuid.UUID), task.Title)
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
	}

}

func GetTasksID(taskStorage *usecase.TaskUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		int_id, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		userID := r.Context().Value("userID")
		res, err := taskStorage.GetTaskByID(userID.(uuid.UUID), int_id)
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
	}

}

func Handlers(r *chi.Mux, taskStorage *usecase.TaskUsecase, userStorage *usecase.AuthUsecase) {
	r.Use(RecoveryMiddleware, LoggingMiddleware)
	r.Get("/ping", Ping())
	r.Post("/echo", Echo())
	r.Get("/hello/{name}", Hello())
	r.Post("/register", Register(userStorage))
	r.Post("/login", Login(userStorage))
	r.Route("/tasks", func(r chi.Router) {
		r.Use(AuthMiddleware)
		r.Get("/", GetTasks(taskStorage))
		r.Post("/", PostTasks(taskStorage))
		r.Get("/{id}", GetTasksID(taskStorage))
	})

}
