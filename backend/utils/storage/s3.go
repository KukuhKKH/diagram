package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Storage struct {
	Client *minio.Client
	Bucket string
	Region string
}

func NewS3Storage(endpoint, accessKey, secretKey, bucket, region string, useSSL bool) (*S3Storage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
		Region: region,
	})

	if err != nil {
		return nil, err
	}

	return &S3Storage{
		Client: client,
		Bucket: bucket,
		Region: region,
	}, nil
}

func (s *S3Storage) Upload(ctx context.Context, filename string, file io.Reader) (string, error) {
	_, err := s.Client.PutObject(context.Background(), s.Bucket, filename, file, -1, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}

	return filename, nil
}

func (s *S3Storage) Delete(filename string) error {
	return s.Client.RemoveObject(context.Background(), s.Bucket, filename, minio.RemoveObjectOptions{})
}

func (s *S3Storage) GetURL(filename string) string {
	return fmt.Sprintf("%s/%s/%s", s.Client.EndpointURL().String(), s.Bucket, filename)
}

func (s *S3Storage) Open(filename string) (io.ReadCloser, error) {
	return s.Client.GetObject(context.Background(), s.Bucket, filename, minio.GetObjectOptions{})
}

func (s *S3Storage) Stat(filename string) (int64, error) {
	info, err := s.Client.StatObject(context.Background(), s.Bucket, filename, minio.StatObjectOptions{})
	if err != nil {
		return 0, err
	}

	return info.Size, nil
}
