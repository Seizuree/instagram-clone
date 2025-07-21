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

func (r *postRepository) CreatePost(post *entities.Post) error {
	return r.db.GetInstance().Create(post).Error
}

func (r *postRepository) GetPostByID(postID uuid.UUID) (*entities.Post, error) {
	var post entities.Post
	if err := r.db.GetInstance().First(&post, "id = ?", postID).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) GetPostsByUserID(userID uuid.UUID) (*[]entities.Post, error) {
	var posts []entities.Post
	if err := r.db.GetInstance().Where("user_id = ?", userID).Find(&posts).Error; err != nil {
		return nil, err
	}
	return &posts, nil
}

func (r *postRepository) UpdatePost(post *entities.Post) error {
	return r.db.GetInstance().Save(post).Error
}

func (r *postRepository) DeletePost(postID uuid.UUID) error {
	return r.db.GetInstance().Delete(&entities.Post{}, "id = ?", postID).Error
}

func (r *postRepository) DeletePostsByUserID(userID uuid.UUID) error {
	return r.db.GetInstance().Where("user_id = ?", userID).Delete(&entities.Post{}).Error
}

func (r *postRepository) CountUserPosts(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.GetInstance().Model(&entities.Post{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}
