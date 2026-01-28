package service

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/parthsarthi-dutt/blog-api/internal/models"
	"github.com/parthsarthi-dutt/blog-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
 // We'll move to .env later

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repository.NewUserRepository(),
	}
}

func (s *AuthService) Register(name, email, password string) error {
	_, err := s.userRepo.FindByEmail(email)
	if err == nil {
		return errors.New("email already in use")
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := models.User{
		Name:      name,
		Email:     email,
		Password:  string(hashed),
		CreatedAt: time.Now().Unix(),
	}

	return s.userRepo.CreateUser(user)
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	return tokenString, err
}
