package postinteraction

import (
	"errors"

	"github.com/LuisFernando-hub/go-social-api/internal/models"
	"github.com/LuisFernando-hub/go-social-api/internal/post"
	"github.com/LuisFernando-hub/go-social-api/internal/types"
)

type Service struct {
	postService *post.Service
	repository  *Repository
}

type InterationPostInput struct {
	Type   string `json:"type"`
	PostID uint   `json:"post_id"`
}

func NewService(repository *Repository, postService *post.Service) *Service {
	return &Service{
		repository:  repository,
		postService: postService,
	}
}

func (s *Service) CreateInteractionPost(input InterationPostInput, userID uint) (*models.Post, error) {
	interactionType := types.InteractionType(input.Type)

	if !interactionType.IsValid() {
		return nil, errors.New("Invalid interaction type")
	}

	post, err := s.postService.GetByID(input.PostID)
	if err != nil {
		return nil, errors.New("Post not found")
	}

	postInteraction, err := s.repository.GetPostInteractionByUserID(userID)

	if err != nil {
		postInteraction = &models.PostInteraction{
			PostID: post.ID,
			UserID: userID,
			Type:   string(types.InteractionType(input.Type)),
		}

		err := s.repository.CreateInteractionPost(postInteraction)
		if err != nil {
			return nil, errors.New("Invalid interaction")
		}

		s.postService.UpdatePostCount(post, types.InteractionType(postInteraction.Type), true)

	} else if types.InteractionType(postInteraction.Type) == interactionType {
		// Usuário clicou novamente na mesma interação.
		// LIKE  -> remove LIKE
		// DISLIKE -> remove DISLIKE

		if err := s.repository.Delete(postInteraction.ID); err != nil {
			return nil, errors.New("invalid interaction")
		}

		s.postService.UpdatePostCount(post, interactionType, false)
	} else {
		// Usuário está trocando a interação.
		// LIKE -> DISLIKE
		// DISLIKE -> LIKE

		oldType := types.InteractionType(postInteraction.Type)

		postInteraction.Type = string(interactionType)

		// Remove uma unidade do post a partir do oldType
		s.postService.UpdatePostCount(post, oldType, false)

		// Adiciona uma unidade do post a partir do oldType
		s.postService.UpdatePostCount(post, interactionType, true)

		
		err := s.repository.UpdateInteractionPost(postInteraction)
		if err != nil {
			return nil, errors.New("Invalid interaction")
		}
	}

	return post, nil
}
