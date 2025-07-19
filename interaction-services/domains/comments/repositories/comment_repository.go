package repositories

import (
	"interaction-services/domains/comments"
	"interaction-services/domains/comments/entities"
	"interaction-services/infrastructures"

	"github.com/google/uuid"
)

type commentRepository struct {
	db infrastructures.Database
}

func NewCommentRepository(db infrastructures.Database) comments.CommentRepository {
	return &commentRepository{db: db}
}

func (c *commentRepository) CreateComment(comment *entities.Comment) (*entities.Comment, error) {
	comment.ID = uuid.New()
	if err := c.db.GetInstance().Create(comment).Error; err != nil {
		return nil, err
	}
	return comment, nil
}

func (c *commentRepository) GetCommentsByPostID(postID uuid.UUID) ([]*entities.Comment, error) {
	var comments []*entities.Comment
	if err := c.db.GetInstance().Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}
