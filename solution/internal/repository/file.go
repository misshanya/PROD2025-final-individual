package repository

import (
	"bytes"
	"context"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
)

type FileRepository struct {
	minioClient *minio.Client
	bucketName  string
}

func NewFileRepository(minioClient *minio.Client, bucketName string) *FileRepository {
	return &FileRepository{
		minioClient: minioClient,
		bucketName:  bucketName,
	}
}

func (r *FileRepository) UploadFile(ctx context.Context, fileKey string, fileContent []byte) error {
	reader := bytes.NewReader(fileContent)
	_, err := r.minioClient.PutObject(ctx, r.bucketName, fileKey, reader, int64(len(fileContent)), minio.PutObjectOptions{})
	return err
}

func (r *FileRepository) GetFileLink(ctx context.Context, fileKey, minioHost string) (string, error) {
	lifetime := time.Hour * 24
	link, err := r.minioClient.PresignedGetObject(ctx, r.bucketName, fileKey, lifetime, url.Values{})
	if err != nil {
		return "", err
	}

	// scheme://userinfo@host/path?query#fragment
	endLink := link.Scheme + "://" + minioHost + link.Path
	return endLink, nil
}
