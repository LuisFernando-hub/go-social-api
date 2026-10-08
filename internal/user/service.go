package user

import (
	"errors"

	"github.com/LuisFernando-hub/go-social-api/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository *Repository
}

type CreateUserInput struct {
	UserName string `json:"username" binding:"required"`
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(input CreateUserInput) (*models.User, error) {
	if input.Email == "" {
		return nil, errors.New("Email is required")
	}

	checkUser, err := s.repository.GetUserByEmail(input.Email)
	if checkUser != nil {
		return nil, errors.New("Email already exists")
	}

	if input.UserName == "" {
		return nil, errors.New("Username is required")
	}

	if input.Password == "" {
		return nil, errors.New("Password is required")
	}

	if len(input.Password) < 6 {
		return nil, errors.New("Password must be 6 characters")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	user := &models.User{
		UserName: input.UserName,
		Email: input.Email,
		PasswordHash: string(hashedPassword),
	}

	if err := s.repository.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GetByID(userID uint) (*models.User, error) {
	if userID == 0 {
		return nil, errors.New("UserID is required")
	}

	user, err := s.repository.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}