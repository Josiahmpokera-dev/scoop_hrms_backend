package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// StorageService handles file storage operations.
// When AWS S3 is configured it delegates to S3; otherwise it uses local disk.
type StorageService struct {
	basePath string
	baseURL  string
	s3       *S3Client
}

// NewStorageService creates a storage service.
// S3 is used automatically when the AWS env vars are set.
func NewStorageService() *StorageService {
	basePath := os.Getenv("STORAGE_PATH")
	if basePath == "" {
		basePath = "./storage"
	}

	baseURL := os.Getenv("STORAGE_URL")
	if baseURL == "" {
		baseURL = "/storage"
	}

	return &StorageService{
		basePath: basePath,
		baseURL:  baseURL,
		s3:       NewS3Client(),
	}
}

// IsS3 returns true when the service is backed by S3.
func (s *StorageService) IsS3() bool {
	return s.s3 != nil
}

// UploadFile uploads a file.
// Returns: fileURL, fileSize, mimeType, error.
// When S3 is configured the returned URL is an absolute S3 URL;
// otherwise it is a relative local path.
func (s *StorageService) UploadFile(file *multipart.FileHeader, folder string, identifier string) (string, int64, string, error) {
	if s.s3 != nil {
		return s.s3.Upload(file, folder, identifier)
	}
	return s.uploadLocal(file, folder, identifier)
}

// DeleteFile deletes a file from storage (S3 or local).
func (s *StorageService) DeleteFile(fileURL string) error {
	if s.s3 != nil && (strings.HasPrefix(fileURL, "https://") || strings.HasPrefix(fileURL, "http://")) {
		return s.s3.Delete(fileURL)
	}
	return s.deleteLocal(fileURL)
}

// ResolveURL converts a stored URL to a full URL suitable for API responses.
// S3 URLs are already absolute and returned as-is.
// Local relative URLs get the request host prepended.
func (s *StorageService) ResolveURL(r *http.Request, storedURL string) string {
	if strings.HasPrefix(storedURL, "https://") || strings.HasPrefix(storedURL, "http://") {
		return storedURL
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	}

	host := r.Host
	if host == "" {
		host = "localhost:8080"
	}

	clean := strings.TrimPrefix(storedURL, "/")
	return fmt.Sprintf("%s://%s/%s", scheme, host, clean)
}

// ---------- Local storage (unchanged legacy behaviour) ----------

func (s *StorageService) uploadLocal(file *multipart.FileHeader, folder string, identifier string) (string, int64, string, error) {
	uploadDir := filepath.Join(s.basePath, folder, identifier)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", 0, "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	ext := filepath.Ext(file.Filename)
	timestamp := time.Now().UnixNano()
	newFilename := fmt.Sprintf("%d%s", timestamp, ext)
	filePath := filepath.Join(uploadDir, newFilename)

	src, err := file.Open()
	if err != nil {
		return "", 0, "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return "", 0, "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	written, err := io.Copy(dst, src)
	if err != nil {
		return "", 0, "", fmt.Errorf("failed to save file: %w", err)
	}

	mimeType := file.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = getMimeTypeFromExtension(ext)
	}

	relativeURL := fmt.Sprintf("%s/%s/%s/%s", s.baseURL, folder, identifier, newFilename)
	return relativeURL, written, mimeType, nil
}

func (s *StorageService) deleteLocal(fileURL string) error {
	filePath := s.urlToPath(fileURL)
	if filePath == "" {
		return fmt.Errorf("invalid file URL")
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func (s *StorageService) urlToPath(fileURL string) string {
	relativePath := strings.TrimPrefix(fileURL, s.baseURL)
	if relativePath == fileURL {
		if idx := strings.Index(fileURL, s.baseURL); idx != -1 {
			relativePath = fileURL[idx+len(s.baseURL):]
		} else {
			return ""
		}
	}
	relativePath = strings.TrimPrefix(relativePath, "/")
	return filepath.Join(s.basePath, relativePath)
}

// GetFilePath returns the full local file path for a given URL.
func (s *StorageService) GetFilePath(fileURL string) string {
	return s.urlToPath(fileURL)
}

// getMimeTypeFromExtension returns the MIME type based on file extension.
func getMimeTypeFromExtension(ext string) string {
	ext = strings.ToLower(ext)
	mimeTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".svg":  "image/svg+xml",
		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".xls":  "application/vnd.ms-excel",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".ppt":  "application/vnd.ms-powerpoint",
		".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
		".txt":  "text/plain",
		".csv":  "text/csv",
		".zip":  "application/zip",
		".rar":  "application/vnd.rar",
		".mp4":  "video/mp4",
		".mp3":  "audio/mpeg",
	}

	if mimeType, ok := mimeTypes[ext]; ok {
		return mimeType
	}
	return "application/octet-stream"
}
