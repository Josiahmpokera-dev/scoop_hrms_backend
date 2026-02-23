package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/config"
	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Client wraps the AWS S3 SDK client.
type S3Client struct {
	client  *s3.Client
	bucket  string
	baseURL string
	prefix  string
}

// NewS3Client builds an S3 client from the application's AWS config.
// Returns nil if S3 is not configured.
func NewS3Client() *S3Client {
	cfg := config.AppConfig
	if cfg == nil || !cfg.AWS.IsS3Enabled() {
		return nil
	}

	awsConf, err := awsCfg.LoadDefaultConfig(context.TODO(),
		awsCfg.WithRegion(cfg.AWS.Region),
		awsCfg.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.AWS.AccessKeyID,
				cfg.AWS.SecretAccessKey,
				"",
			),
		),
	)
	if err != nil {
		fmt.Printf("[S3] WARNING: failed to load AWS config: %v — falling back to local storage\n", err)
		return nil
	}

	baseURL := cfg.AWS.S3BaseURL
	if baseURL == "" {
		baseURL = fmt.Sprintf("https://%s.s3.%s.amazonaws.com", cfg.AWS.S3Bucket, cfg.AWS.Region)
	}

	prefix := strings.TrimSuffix(cfg.AWS.S3DocsPrefix, "/")
	if prefix != "" {
		prefix += "/"
	}

	return &S3Client{
		client:  s3.NewFromConfig(awsConf),
		bucket:  cfg.AWS.S3Bucket,
		baseURL: strings.TrimSuffix(baseURL, "/"),
		prefix:  prefix,
	}
}

// Upload puts a multipart file on S3 and returns the public URL, size, and MIME type.
func (s *S3Client) Upload(file *multipart.FileHeader, folder, identifier string) (string, int64, string, error) {
	src, err := file.Open()
	if err != nil {
		return "", 0, "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	ext := filepath.Ext(file.Filename)
	timestamp := time.Now().UnixNano()
	key := fmt.Sprintf("%s%s/%s/%d%s", s.prefix, folder, identifier, timestamp, ext)

	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = getMimeTypeFromExtension(ext)
	}

	_, err = s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      &s.bucket,
		Key:         &key,
		Body:        src,
		ContentType: &contentType,
	})
	if err != nil {
		return "", 0, "", fmt.Errorf("S3 upload failed: %w", err)
	}

	url := fmt.Sprintf("%s/%s", s.baseURL, key)
	return url, file.Size, contentType, nil
}

// UploadReader puts an io.Reader on S3 (used for programmatic uploads, not multipart).
func (s *S3Client) UploadReader(body io.Reader, key, contentType string) (string, error) {
	fullKey := s.prefix + key

	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      &s.bucket,
		Key:         &fullKey,
		Body:        body,
		ContentType: &contentType,
	})
	if err != nil {
		return "", fmt.Errorf("S3 upload failed: %w", err)
	}

	return fmt.Sprintf("%s/%s", s.baseURL, fullKey), nil
}

// Delete removes a file from S3 given its full URL.
func (s *S3Client) Delete(fileURL string) error {
	key := strings.TrimPrefix(fileURL, s.baseURL+"/")
	if key == fileURL {
		return fmt.Errorf("URL does not belong to this S3 bucket")
	}

	_, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	})
	return err
}
