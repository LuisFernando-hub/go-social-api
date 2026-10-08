package post

import (
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

func (r *Repository) GetByID(postID uint) (*models.Post, error) {
	var post *models.Post
	if err := r.db.Where("id = ?", postID).First(&post).Error; err != nil {
		return nil, err
	}

	return post, nil
}

func (r *Repository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

func (r *Repository) Update(post *models.Post) error {
	return r.db.Save(&post).Error
}

func (r *Repository) List() ([]*models.Post, error) {
	var posts []*models.Post
	
	if err := r.db.Preload("Interactions").Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}