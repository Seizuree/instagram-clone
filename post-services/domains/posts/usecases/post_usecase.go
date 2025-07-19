package usecases

import (
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"post-services/domains/posts"
	"post-services/domains/posts/entities"
	"post-services/infrastructures"
	"strings"

	"github.com/google/uuid"
)

type postUseCase struct {
	postRepo    posts.PostRepository
	minioClient *infrastructures.MinioClient
}

func NewPostUseCase(postRepo posts.PostRepository, minioClient *infrastructures.MinioClient) posts.PostUseCase {
	return &postUseCase{postRepo: postRepo, minioClient: minioClient}
}

// CreatePost implements posts.PostUseCase.
func (p *postUseCase) CreatePost(userID uuid.UUID, caption string, fileHeader *multipart.FileHeader) (*entities.Post, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	objectName := fmt.Sprintf("%s-%s", uuid.New().String(), fileHeader.Filename)

	imageURL, err := p.minioClient.UploadFile(objectName, file, fileHeader.Size)
	if err != nil {
		return nil, err
	}

	post := &entities.Post{
		UserID:   userID,
		ImageURL: imageURL,
		Caption:  caption,
	}

	if err := p.postRepo.CreatePost(post); err != nil {
		return nil, err
	}

	return post, nil
}

// GetPost implements posts.PostUseCase.
func (p *postUseCase) GetPost(postID uuid.UUID) (*entities.Post, error) {
	return p.postRepo.GetPostByID(postID)
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
	log.Printf("Attempting to delete all posts for user %s", userID)

	// First, get all posts for the user to delete associated files from Minio.
	posts, err := p.postRepo.GetPostsByUserID(userID)
	if err != nil {
		log.Printf("Error getting posts for user %s: %v", userID, err)
		return fmt.Errorf("could not get posts for user %s: %w", userID, err)
	}

	// Delete all associated files from Minio.
	for _, post := range *posts {
		urlParts := strings.Split(post.ImageURL, "/")
		if len(urlParts) > 0 {
			objectName := urlParts[len(urlParts)-1]
			if err := p.minioClient.DeleteFile(objectName); err != nil {
				// Log the error and continue. The database record will be deleted anyway.
				log.Printf("ERROR: failed to delete file '%s' from Minio for post %s: %v", objectName, post.ID, err)
			}
		} else {
			log.Printf("WARN: could not parse object name from URL: %s for post %s", post.ImageURL, post.ID)
		}
	}

	// After attempting to delete all files, delete all post records from the database.
	if err := p.postRepo.DeletePostsByUserID(userID); err != nil {
		log.Printf("Error deleting posts from repository for user %s: %v", userID, err)
		return fmt.Errorf("could not delete posts from database for user %s: %w", userID, err)
	}

	log.Printf("Successfully deleted all posts and associated files for user %s", userID)
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

	// Extract object name from URL
	urlParts := strings.Split(post.ImageURL, "/")
	if len(urlParts) == 0 {
		log.Printf("WARN: could not parse object name from URL: %s", post.ImageURL)
	} else {
		objectName := urlParts[len(urlParts)-1]
		if err := p.minioClient.DeleteFile(objectName); err != nil {
			// Log the error but proceed with DB deletion.
			// In a production system, you might add this to a retry queue.
			log.Printf("ERROR: failed to delete file '%s' from Minio: %v", objectName, err)
		}
	}

	return p.postRepo.DeletePost(postID)
}
