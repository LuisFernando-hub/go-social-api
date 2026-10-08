package post

import (
	"errors"

	"github.com/LuisFernando-hub/go-social-api/internal/models"
	"github.com/LuisFernando-hub/go-social-api/internal/types"
	"github.com/LuisFernando-hub/go-social-api/internal/user"
)

type Service struct {
	respository *Repository
	userService *user.Service
}

type CreatePostInput struct {
	Content string `json:"content" binding:"required"`
}

func NewService(repository *Repository, userService *user.Service) *Service {
	return &Service{
		respository: repository,
		userService: userService,
	}
}

func (s *Service) GetByID(userID uint) (*models.Post, error) {
	if userID == 0 {
		return nil, errors.New("user_id is required")
	}

	var post *models.Post

	post, err := s.respository.GetByID(userID)

	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *Service) Create(input CreatePostInput, userID uint) (*models.Post, error) {
	if userID == 0 {
		return nil, errors.New("user_id is required")
	}

	if input.Content == "" {
		return nil, errors.New("content is required")
	}

	if len(input.Content) > 5000 {
		return nil, errors.New("content is too long")
	}

	user, err := s.userService.GetByID(userID)

	if err != nil {
		return nil, errors.New("user not found")
	}

	post := &models.Post{
		UserID:  user.ID,
		Content: input.Content,
	}

	if err := s.respository.Create(post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *Service) List() ([]*models.Post, error) {
	posts, err := s.respository.List()

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *Service) UpdatePostCount(post *models.Post, interactionType types.InteractionType, plus bool) (*models.Post, error) {
	if interactionType == types.InteractionLike {
		if plus {
            post.LikeCount++
        } else {
            post.LikeCount--
        }
	} 

	if interactionType == types.InteractionDeslike {

        if plus {
            post.DeslikeCount++
        } else {
            post.DeslikeCount--
        }
    }

	if err := s.respository.Update(post); err != nil {
		return nil, err
	}

	return post, nil
}
