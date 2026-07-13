package domain

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var UserExists = errors.New("User already exists")
var UserNotExists = errors.New("User not exists")
var ErrTaskNotFound = errors.New("Task not found")

type Answer struct {
	Message string `json:"message,omitempty"`
	Text    string `json:"text,omitempty"`
	Hello   string `json:"hello,omitempty"`
}

type Task struct {
	UserId    uuid.UUID `json:"user_id,omitempty"`
	ID        int       `json:"id,omitempty"`
	Title     string    `json:"title,omitempty"`
	Done      bool      `json:"done,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type User struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Username  string    `json:"username,omitempty"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type CustomClaims struct {
	Userid   uuid.UUID
	Username string
	jwt.RegisteredClaims
}
