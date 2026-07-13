package usecase

import (
	"errors"
	"log"
	"os"
	"server/domain"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var BadRequest = errors.New("Bad request")
var InternalServerError = errors.New("Internal Server Error")
var ErrInvalidCredentials = errors.New("Err Invalid Credentials")

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
		return nil, BadRequest
	}
	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, InternalServerError
	}
	user := &domain.User{
		ID:        uuid.New(),
		Username:  username,
		Password:  string(HashedPassword),
		CreatedAt: time.Now(),
	}
	_, err = uu.userrepo.SaveUser(user)
	if err != nil {
		return nil, InternalServerError
	}

	return user, nil
}

func (uu *AuthUsecase) Login(username, password string) (*string, error) {
	if username == "" || len(username) < 3 {
		return nil, BadRequest
	}
	if len(password) > 60 || len(password) < 6 {
		return nil, BadRequest
	}
	user, err := uu.userrepo.FindByUsername(username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	claims := &domain.CustomClaims{
		Userid:   user.ID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(os.Getenv("SECRETKEY")))
	if err != nil {
		log.Printf("Ошибка подписи токена: %v", err)
		return nil, InternalServerError
	}
	return &ss, nil
}
