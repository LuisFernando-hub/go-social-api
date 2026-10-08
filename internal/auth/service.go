package auth

import (
	"errors"
	"log"
	"time"

	"github.com/LuisFernando-hub/go-social-api/internal/config"
	"github.com/LuisFernando-hub/go-social-api/internal/user"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)


type Service struct {
	repository *user.Repository
}

type LoginAuthInput struct {
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID uint `json:"id"`
	UserName string `json:"username"`
	Email string `json:"email"`
}

type AuthResponse struct {
	User *UserResponse `json:"user"`
	Token string `json:"token"`
}

func NewService(repository *user.Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Login(cfg *config.Config, input LoginAuthInput) (*AuthResponse, error) {
	if input.Email == "" {
		return nil, errors.New("Email is required")
	}

	if input.Password == "" {
		return nil, errors.New("Password is required")
	}

	user, err := s.repository.GetUserByEmail(input.Email)
	if err != nil {
		return nil, err
	}

	log.Printf("USER: %v", user);
	log.Printf("INPUT: %v", input);

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return nil, err
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email": user.Email,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: tokenString,
		User: &UserResponse{
			ID: user.ID,
			UserName: user.Email,
			Email: user.Email,
		},
	}, nil
}