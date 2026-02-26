package storage

import (
	"context"
	"fmt"
	"io"

	"git.dev.siap.id/kukuhkkh/app-diagram/utils/config"
)

type Storage interface {
	Upload(ctx context.Context, filename string, file io.Reader) (string, error)
	Delete(filename string) error
	GetURL(filename string) string
	Open(filename string) (io.ReadCloser, error)
	Stat(filename string) (int64, error)
}

func NewStorage(cfg *config.Config) (Storage, error) {
	switch cfg.Storage.Driver {
	case "local":
		return NewLocalStorage(cfg.Storage.Local.Path)
	case "ftp":
		return NewSftpStorage(
			cfg.Storage.Ftp.Host,
			cfg.Storage.Ftp.Port,
			cfg.Storage.Ftp.User,
			cfg.Storage.Ftp.Password,
			cfg.Storage.Ftp.BaseDir,
			cfg.Storage.Ftp.PublicUrl,
		), nil
	case "s3":
		return NewS3Storage(
			cfg.Storage.S3.Endpoint,
			cfg.Storage.S3.AccessKey,
			cfg.Storage.S3.SecretKey,
			cfg.Storage.S3.Bucket,
			cfg.Storage.S3.Region,
			cfg.Storage.S3.UseSsl,
		)
	default:
		return nil, fmt.Errorf("storage driver %s not supported", cfg.Storage.Driver)
	}
}
