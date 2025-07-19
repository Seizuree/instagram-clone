package repositories

import (
	"post-services/domains/posts"
	"post-services/domains/posts/entities"
	"post-services/infrastructures"

	"github.com/google/uuid"
)

type postRepository struct {
	db infrastructures.Database
}

func NewPostRepository(db infrastructures.Database) posts.PostRepository {
	return &postRepository{db: db}
}

// CreatePost implements posts.PostRepository.
func (p *postRepository) CreatePost(post *entities.Post) error {
	return p.db.GetInstance().Create(post).Error
}

// GetPostByID implements posts.PostRepository.
func (p *postRepository) GetPostByID(postID uuid.UUID) (*entities.Post, error) {
	var post entities.Post
	if err := p.db.GetInstance().First(&post, "id = ?", postID).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

// GetPostsByUserID implements posts.PostRepository.
func (p *postRepository) GetPostsByUserID(userID uuid.UUID) (*[]entities.Post, error) {
	var posts []entities.Post
	if err := p.db.GetInstance().Where("user_id = ?", userID).Find(&posts).Error; err != nil {
		return nil, err
	}
	return &posts, nil
}

// UpdatePost implements posts.PostRepository.
func (p *postRepository) UpdatePost(post *entities.Post) error {
	return p.db.GetInstance().Save(post).Error
}

// DeletePostsByUserID implements posts.PostRepository.
func (p *postRepository) DeletePostsByUserID(userID uuid.UUID) error {
	// This performs a bulk delete of all posts matching the user_id.
	return p.db.GetInstance().Where("user_id = ?", userID).Delete(&entities.Post{}).Error
}

// DeletePost implements posts.PostRepository.
func (p *postRepository) DeletePost(postID uuid.UUID) error {
	return p.db.GetInstance().Delete(&entities.Post{}, "id = ?", postID).Error
}
