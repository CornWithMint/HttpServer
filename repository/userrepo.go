package repository

import (
	"errors"
	"server/domain"
	"sync"

	"github.com/google/uuid"
)

var UserExists = errors.New("User already exists")
var UserNotExists = errors.New("User not exists")
var BadRequest = errors.New("Bad request")

type UserStorage struct {
	users map[uuid.UUID]*domain.User
	mu    *sync.RWMutex
}

func NewUserStorage() *UserStorage {
	return &UserStorage{
		users: make(map[uuid.UUID]*domain.User),
		mu:    &sync.RWMutex{},
	}
}

func (us *UserStorage) SaveUser(user *domain.User) (*domain.User, error) {
	us.mu.Lock()
	defer us.mu.Unlock()

	us.users[user.ID] = user
	return user, nil
}

func (us *UserStorage) FindByUsername(username string) (*domain.User, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()
	for _, v := range us.users {
		if v.Username == username {
			return v, nil
		}
	}
	return nil, UserNotExists
}
