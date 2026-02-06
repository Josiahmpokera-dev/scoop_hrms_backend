package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// StorageService handles file storage operations
type StorageService struct {
	basePath string
	baseURL  string
}

// NewStorageService creates a new storage service
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
	}
}

// UploadFile uploads a file to storage
// Returns: fileURL, fileSize, mimeType, error
func (s *StorageService) UploadFile(file *multipart.FileHeader, folder string, identifier string) (string, int64, string, error) {
	// Create directory structure
	uploadDir := filepath.Join(s.basePath, folder, identifier)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", 0, "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	timestamp := time.Now().UnixNano()
	newFilename := fmt.Sprintf("%d%s", timestamp, ext)
	filePath := filepath.Join(uploadDir, newFilename)

	// Open source file
	src, err := file.Open()
	if err != nil {
		return "", 0, "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", 0, "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy file contents
	written, err := io.Copy(dst, src)
	if err != nil {
		return "", 0, "", fmt.Errorf("failed to save file: %w", err)
	}

	// Get mime type from file header or extension
	mimeType := file.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = getMimeTypeFromExtension(ext)
	}

	// Generate relative URL
	relativeURL := fmt.Sprintf("%s/%s/%s/%s", s.baseURL, folder, identifier, newFilename)

	return relativeURL, written, mimeType, nil
}

// DeleteFile deletes a file from storage
func (s *StorageService) DeleteFile(fileURL string) error {
	// Convert URL to file path
	filePath := s.urlToPath(fileURL)
	if filePath == "" {
		return fmt.Errorf("invalid file URL")
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // File doesn't exist, consider it deleted
	}

	// Delete the file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// urlToPath converts a storage URL to a file path
func (s *StorageService) urlToPath(fileURL string) string {
	// Remove base URL prefix to get relative path
	relativePath := strings.TrimPrefix(fileURL, s.baseURL)
	if relativePath == fileURL {
		// URL might be a full URL, try to extract path
		if idx := strings.Index(fileURL, s.baseURL); idx != -1 {
			relativePath = fileURL[idx+len(s.baseURL):]
		} else {
			return ""
		}
	}

	// Clean the path and join with base path
	relativePath = strings.TrimPrefix(relativePath, "/")
	return filepath.Join(s.basePath, relativePath)
}

// GetFilePath returns the full file path for a given URL
func (s *StorageService) GetFilePath(fileURL string) string {
	return s.urlToPath(fileURL)
}

// getMimeTypeFromExtension returns the MIME type based on file extension
func getMimeTypeFromExtension(ext string) string {
	ext = strings.ToLower(ext)
	mimeTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".xls":  "application/vnd.ms-excel",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".txt":  "text/plain",
		".csv":  "text/csv",
	}

	if mimeType, ok := mimeTypes[ext]; ok {
		return mimeType
	}
	return "application/octet-stream"
}
