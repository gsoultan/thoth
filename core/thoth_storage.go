package core

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// ThothStorage implements the StorageProvider interface for multiple storage types.
type ThothStorage struct {
	s3Config *S3Config
}

// NewThothStorage creates a new ThothStorage with the provided S3 configuration.
func NewThothStorage(s3Config *S3Config) *ThothStorage {
	return &ThothStorage{s3Config: s3Config}
}

// Open opens a reader for the given URI, supporting local files, HTTP, and S3.
func (s *ThothStorage) Open(ctx context.Context, uri string) (io.ReadCloser, error) {
	u, err := url.Parse(uri)
	if err != nil || u.Scheme == "" {
		// Treat as local file
		return os.Open(uri)
	}

	switch u.Scheme {
	case "http", "https":
		return s.openHTTP(ctx, uri)
	case "s3":
		if s.s3Config == nil {
			return nil, fmt.Errorf("s3 config not provided for s3 uri: %s", uri)
		}
		return s.openS3(ctx, u)
	default:
		// Fallback to local file for unknown schemes
		return os.Open(uri)
	}
}

// Save writes the content of reader to the given URI.
func (s *ThothStorage) Save(ctx context.Context, uri string, reader io.Reader) error {
	u, err := url.Parse(uri)
	if err != nil || u.Scheme == "" {
		// Treat as local file
		return s.saveLocal(uri, reader)
	}

	switch u.Scheme {
	case "s3":
		if s.s3Config == nil {
			return fmt.Errorf("s3 config not provided for s3 uri: %s", uri)
		}
		return s.saveS3(ctx, u, reader)
	case "http", "https":
		return fmt.Errorf("saving to http/https not supported")
	default:
		// Fallback to local file for unknown schemes
		return s.saveLocal(uri, reader)
	}
}

func (s *ThothStorage) saveLocal(uri string, reader io.Reader) error {
	f, err := os.Create(uri)
	if err != nil {
		return fmt.Errorf("create local file: %w", err)
	}
	defer f.Close()

	_, err = io.Copy(f, reader)
	if err != nil {
		return fmt.Errorf("copy local file: %w", err)
	}
	return nil
}

func (s *ThothStorage) openHTTP(ctx context.Context, uri string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("create http request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute http request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("http error status: %s", resp.Status)
	}
	return resp.Body, nil
}

func (s *ThothStorage) openS3(ctx context.Context, u *url.URL) (io.ReadCloser, error) {
	client, err := minio.New(s.s3Config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s.s3Config.AccessKeyID, s.s3Config.SecretAccessKey, ""),
		Secure: s.s3Config.UseSSL,
		Region: s.s3Config.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("create minio client: %w", err)
	}

	bucket := u.Host
	key := strings.TrimPrefix(u.Path, "/")

	obj, err := client.GetObject(ctx, bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("get s3 object: %w", err)
	}

	_, err = obj.Stat()
	if err != nil {
		obj.Close()
		return nil, fmt.Errorf("stat s3 object: %w", err)
	}

	return obj, nil
}

func (s *ThothStorage) saveS3(ctx context.Context, u *url.URL, reader io.Reader) error {
	client, err := minio.New(s.s3Config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s.s3Config.AccessKeyID, s.s3Config.SecretAccessKey, ""),
		Secure: s.s3Config.UseSSL,
		Region: s.s3Config.Region,
	})
	if err != nil {
		return fmt.Errorf("create minio client: %w", err)
	}

	bucket := u.Host
	key := strings.TrimPrefix(u.Path, "/")

	_, err = client.PutObject(ctx, bucket, key, reader, -1, minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("put s3 object: %w", err)
	}

	return nil
}
