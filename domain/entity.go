package domain

import "time"

type Answer struct {
	Message string `json:"message,omitempty"`
	Text    string `json:"text,omitempty"`
	Hello   string `json:"hello,omitempty"`
}
type Task struct {
	UserId    int       `json:"user_id,omitempty"`
	ID        int       `json:"id,omitempty"`
	Title     string    `json:"title,omitempty"`
	Done      bool      `json:"done,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type User struct {
	ID        int       `json:"id,omitempty"`
	Username  string    `json:"username,omitempty"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
