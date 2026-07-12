package usecase

import (
	"errors"
	"server/domain"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var UserExists = errors.New("User already exists")
var UserNotExists = errors.New("User not exists")
var BadRequest = errors.New("Bad request")
var InternalServerError = errors.New("Internal Server Error")

type UserRepository interface {
	SaveUser(user *domain.User) (*domain.User, error)
	FindByUsername(username string) (*domain.User, error)
}

type AuthUsecase struct {
	userrepo UserRepository
	mu       *sync.RWMutex
}

func NewAuthUsecase(ur UserRepository) *AuthUsecase {
	return &AuthUsecase{
		userrepo: ur,
		mu:       &sync.RWMutex{},
	}
}

func (uu *AuthUsecase) Register(username, password string) (*domain.User, error) {
	if username == "" || len(username) < 3 {
		return nil, BadRequest
	}
	if len(password) > 60 || len(password) < 6 {
		return nil, BadRequest
	}

	exist, _ := uu.userrepo.FindByUsername(username)
	if exist != nil {
		return nil, UserExists
	}
	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, InternalServerError
	}
	user := &domain.User{
		ID:        1,
		Username:  username,
		Password:  string(HashedPassword),
		CreatedAt: time.Now(),
	}
	us, err := uu.userrepo.SaveUser(user)
	if err != nil {
		return nil, InternalServerError
	}

}
