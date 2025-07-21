package http

import (
	"errors"
	"net/http"
	"post-services/domains/posts"
	"post-services/domains/posts/models/requests"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PostHttp struct {
	pc posts.PostUseCase
}

func NewPostHttp(pc posts.PostUseCase) *PostHttp {
	return &PostHttp{pc}
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

func (h *PostHttp) CreatePost(c *gin.Context) {
	userID, err := getUserIDFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	fileHeader, err := c.FormFile("image")
	caption := c.PostForm("caption")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image file is required"})
		return
	}

	const maxFileSize = 10 * 1024 * 1024
	if fileHeader.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file size exceeds the 10MB limit"})
		return
	}

	post, err := h.pc.CreatePost(userID, caption, fileHeader)
	if err != nil {
		// Check for specific error types if the use case provides them
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, post)
}

func (h *PostHttp) GetPost(c *gin.Context) {
	_, err := getUserIDFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	postID, err := uuid.Parse(c.Param("post_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	post, err := h.pc.GetPost(postID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *PostHttp) GetPostsByUser(c *gin.Context) {
	userID, err := getUserIDFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	posts, err := h.pc.GetPostsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *PostHttp) UpdatePost(c *gin.Context) {
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

	var req requests.PostUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post, err := h.pc.UpdatePost(userID, postID, req.Caption)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *PostHttp) DeletePost(c *gin.Context) {
	postID, err := uuid.Parse(c.Param("post_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	userID, err := getUserIDFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := h.pc.DeletePost(userID, postID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post deleted successfully"})
}

func (h *PostHttp) CountUserPosts(c *gin.Context) {
	_, err := getUserIDFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	count, err := h.pc.CountUserPosts(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count user posts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}
