package postinteraction

import (
	"log"

	"github.com/LuisFernando-hub/go-social-api/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateInteractionPost(postInteraction *models.PostInteraction) error {
	return r.db.Create(postInteraction).Error
}

func (r *Repository) UpdateInteractionPost(postInteraction *models.PostInteraction) error {
	return r.db.Save(postInteraction).Error
}

func (r *Repository) GetPostInteractionByUserID(userID uint) (*models.PostInteraction, error) {
	var postInteraction *models.PostInteraction
	if err := r.db.Where("user_id = ?", userID).First(&postInteraction).Error; err != nil {
		return nil, err
	}

	return postInteraction, nil
}

func (r *Repository) GetByID(postInterationID uint) (*models.PostInteraction, error) {
	var postInteraction *models.PostInteraction

	if err := r.db.Where("id = ?", postInterationID).First(&postInteraction).Error; err != nil {
		return nil, err
	}

	return postInteraction, nil
}

func (r *Repository) Delete(postInteractionID uint) error {
	var postInteraction *models.PostInteraction

	log.Println("postInteractionID: ", postInteractionID)

	postInteraction, err := r.GetByID(postInteractionID)
	if err != nil {
		return err
	}

	return r.db.Delete(&postInteraction).Error
}
