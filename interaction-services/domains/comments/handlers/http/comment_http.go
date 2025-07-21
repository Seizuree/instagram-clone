package http

import (
	"errors"
	"interaction-services/domains/comments"
	"interaction-services/domains/comments/models/requests"
	"strings"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CommentHttp struct {
	cc comments.CommentUseCase
}

func NewCommentHttp(cc comments.CommentUseCase) *CommentHttp {
	return &CommentHttp{cc}
}

func getUserIDFromHeader(c *gin.Context) (uuid.UUID, error) {
	userIDStr := c.GetHeader("X-User-ID")
	if userIDStr == "" {
		return uuid.Nil, errors.New("X-User-ID header is missing, request must come from the API gateway")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, errors.New("invalid user id format in X-User-ID header")
	}
	return userID, nil
}

func (h *CommentHttp) CreateComment(c *gin.Context) {
	userID, err := getUserIDFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	postID, err := uuid.Parse(c.Param("post_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	var req requests.CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment, err := h.cc.CreateComment(userID, postID, req.Comment)
	if err != nil {
		if strings.Contains(err.Error(), "post not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func (h *CommentHttp) GetCommentsByPostID(c *gin.Context) {
	postID, err := uuid.Parse(c.Param("post_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	comments, err := h.cc.GetCommentsByPostID(postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}
