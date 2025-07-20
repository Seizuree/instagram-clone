package infrastructures

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"post-services/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	client     *minio.Client
	bucketName string
}

func NewMinioClient(cfg *config.Config) *MinioClient {
	minioClient, err := minio.New(cfg.Minio.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Minio.MinioAccessKey, cfg.Minio.MinioSecretKey, ""),
		Secure: cfg.Minio.MinioUseSSL,
	})

	if err != nil {
		log.Fatalf("failed to initialize minio client: %v", err)
	}

	// Create bucket if it doesn't exist
	exists, err := minioClient.BucketExists(context.Background(), cfg.Minio.MinioBucketName)
	if err != nil {
		log.Fatalf("failed to check if bucket exists: %v", err)
	}
	if !exists {
		err = minioClient.MakeBucket(context.Background(), cfg.Minio.MinioBucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("failed to create bucket: %v", err)
		}
	}

	return &MinioClient{
		client:     minioClient,
		bucketName: cfg.Minio.MinioBucketName,
	}
}

func (m *MinioClient) UploadFile(objectName string, file multipart.File, fileSize int64) (string, error) {
	_, err := m.client.PutObject(context.Background(), m.bucketName, objectName, file, fileSize, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/%s", m.client.EndpointURL(), m.bucketName, objectName), nil
}

func (m *MinioClient) DeleteFile(objectName string) error {
	opts := minio.RemoveObjectOptions{
		GovernanceBypass: true,
	}
	return m.client.RemoveObject(context.Background(), m.bucketName, objectName, opts)
}

func (m *MinioClient) UploadBytes(objectName string, data []byte) (string, error) {
    reader := bytes.NewReader(data)

    _, err := m.client.PutObject(context.Background(), m.bucketName, objectName, reader, int64(len(data)), minio.PutObjectOptions{ContentType: "image/jpeg"})
    if err != nil {
        return "", err
    }

    return m.endpoint + "/" + m.bucketName + "/" + objectName, nil
}
