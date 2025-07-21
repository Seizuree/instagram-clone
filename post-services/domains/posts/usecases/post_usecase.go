package usecases

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"post-services/config"
	"post-services/domains/posts"
	"post-services/domains/posts/entities"
	"post-services/domains/posts/models/responses"
	"post-services/infrastructures"
	"strings"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

type postUseCase struct {
	postRepo    posts.PostRepository
	minioClient *infrastructures.MinioClient
	rabbitMQ    *infrastructures.RabbitMQ
	config      *config.Config
}

func NewPostUseCase(postRepo posts.PostRepository, minioClient *infrastructures.MinioClient, rabbitMQ *infrastructures.RabbitMQ, config *config.Config) posts.PostUseCase {
	return &postUseCase{postRepo: postRepo, minioClient: minioClient, rabbitMQ: rabbitMQ, config: config}
}

// CreatePost implements posts.PostUseCase.
func (p *postUseCase) CreatePost(userID uuid.UUID, caption string, fileHeader *multipart.FileHeader) (*entities.Post, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	postID := uuid.New()

	fileNameParts := strings.Split(fileHeader.Filename, ".")
	fileExtension := ""

	if len(fileNameParts) > 1 {
		fileExtension = "." + fileNameParts[len(fileNameParts)-1]
	}

	objectName := fmt.Sprintf("%s/%s/main%s", userID.String(), postID.String(), fileExtension)
	thumbObjectName := fmt.Sprintf("%s/%s/thumb%s", userID.String(), postID.String(), fileExtension)

	imageURL, err := p.minioClient.UploadFile(objectName, file, fileHeader.Size)
	if err != nil {
		return nil, err
	}

	thumbFile, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer thumbFile.Close()

	thumbnailBytes, err := generateThumbnail(thumbFile)
	if err != nil {
		return nil, err
	}

	thumbURL, err := p.minioClient.UploadBytes(thumbObjectName, thumbnailBytes)
	if err != nil {
		return nil, err
	}

	post := &entities.Post{
		ID:       postID,
		UserID:   userID,
		ImageURL: imageURL,
		ThumbURL: thumbURL,
		Caption:  caption,
	}

	if err := p.postRepo.CreatePost(post); err != nil {
		return nil, err
	}

	queueName := "post.created"
	message := map[string]interface{}{
		"post_id":    post.ID.String(),
		"user_id":    post.UserID.String(),
		"image_url":  post.ImageURL,
		"thumb_url":  post.ThumbURL,
		"caption":    post.Caption,
		"created_at": post.CreatedAt,
	}
	if err := p.rabbitMQ.PublishJSON(context.Background(), queueName, message); err != nil {
		log.Printf("CRITICAL: Failed to publish post.created event for postID %s: %v", post.ID, err)
	}

	return post, nil
}

// GetPost implements posts.PostUseCase.
func (p *postUseCase) GetPost(postID uuid.UUID) (*responses.PostDetailResponse, error) {
	post, err := p.postRepo.GetPostByID(postID)
	if err != nil {
		return nil, err
	}

	likeCount, commentCount, err := p.getInteractionCounts(postID)
	if err != nil {
		log.Printf("WARN: Failed to get interaction counts for post %s: %v", postID, err)
		likeCount = 0
		commentCount = 0
	}

	return &responses.PostDetailResponse{
		ID:           post.ID.String(),
		UserID:       post.UserID.String(),
		ImageURL:     post.ImageURL,
		Caption:      post.Caption,
		LikeCount:    likeCount,
		CommentCount: commentCount,
	}, nil
}

// GetPostsByUserID implements posts.PostUseCase.
func (p *postUseCase) GetPostsByUserID(userID uuid.UUID) (*[]entities.Post, error) {
	return p.postRepo.GetPostsByUserID(userID)
}

// UpdatePost implements posts.PostUseCase.
func (p *postUseCase) UpdatePost(userID uuid.UUID, postID uuid.UUID, caption string) (*entities.Post, error) {
	post, err := p.postRepo.GetPostByID(postID)

	if err != nil {
		return nil, err
	}

	if post.UserID != userID {
		return nil, errors.New("user not authorized to update this post")
	}

	post.Caption = caption
	if err := p.postRepo.UpdatePost(post); err != nil {
		return nil, err
	}

	return post, nil
}

// DeletePostsByUserID implements posts.PostUseCase.
func (p *postUseCase) DeletePostsByUserID(userID uuid.UUID) error {
	userFolderPrefix := userID.String() + "/"

	if err := p.minioClient.DeleteUserFolder(userFolderPrefix); err != nil {
		// Log the error but continue, as we still want to clear the DB records.
		log.Printf("ERROR: Failed during bulk file deletion for user %s: %v", userID, err)
	}

	if err := p.postRepo.DeletePostsByUserID(userID); err != nil {
		log.Printf("Error deleting posts from repository for user %s: %v", userID, err)
		return fmt.Errorf("could not delete posts from database for user %s: %w", userID, err)
	}

	log.Printf("Successfully processed deletion for user %s", userID)
	return nil
}

// DeletePost implements posts.PostUseCase.
func (p *postUseCase) DeletePost(userID uuid.UUID, postID uuid.UUID) error {
	post, err := p.postRepo.GetPostByID(postID)
	if err != nil {
		return err
	}

	if post.UserID != userID {
		return errors.New("user not authorized to delete this post")
	}

	extractObjectKey := func(fileURL string) string {
		parsedURL, err := url.Parse(fileURL)
		if err != nil {
			return ""
		}
		parts := strings.Split(strings.TrimPrefix(parsedURL.Path, "/"), "/")
		if len(parts) > 1 {
			return strings.Join(parts[1:], "/")
		}
		return ""
	}

	objectName := extractObjectKey(post.ImageURL)
	thumbName := extractObjectKey(post.ThumbURL)

	if objectName != "" {
		if err := p.minioClient.DeleteFile(objectName); err != nil {
			log.Printf("ERROR: Failed to delete main image '%s' from Minio for post %s: %v", objectName, postID, err)
		}
	} else {
		log.Printf("WARN: Could not extract object key from ImageURL: %s", post.ImageURL)
	}

	if thumbName != "" {
		if err := p.minioClient.DeleteFile(thumbName); err != nil {
			log.Printf("ERROR: Failed to delete thumbnail '%s' from Minio for post %s: %v", thumbName, postID, err)
		}
	} else {
		log.Printf("WARN: Could not extract object key from ThumbURL: %s", post.ThumbURL)
	}

	return p.postRepo.DeletePost(postID)
}

func generateThumbnail(file multipart.File) ([]byte, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	thumbnail := resize.Resize(300, 0, img, resize.Lanczos3)

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, thumbnail, nil); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (p *postUseCase) CountUserPosts(userID uuid.UUID) (int64, error) {
	return p.postRepo.CountUserPosts(userID)
}

func (p *postUseCase) getInteractionCounts(postID uuid.UUID) (likeCount int64, commentCount int64, err error) {
	interactionURL := p.config.Server.InteractionServiceURL
	resp, err := http.Get(fmt.Sprintf("%s/api/internal/interactions/%s/counts", interactionURL, postID))
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("interaction-service returned non-OK status: %d", resp.StatusCode)
	}

	var result responses.PostInteractionCounts

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, 0, err
	}

	return result.LikeCount, result.CommentCount, nil
}
