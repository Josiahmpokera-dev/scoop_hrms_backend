package handlers

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/response"
	"github.com/Josiahmpokera-dev/hrms-backend/internal/utils/storage"
	"github.com/gin-gonic/gin"
)

var allowedImageExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true, ".svg": true,
}

var allowedDocExts = map[string]bool{
	".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true,
	".ppt": true, ".pptx": true, ".txt": true, ".csv": true, ".zip": true, ".rar": true,
}

const maxFileSize = 25 << 20 // 25 MB

type UploadHandler struct {
	storage *storage.StorageService
}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{storage: storage.NewStorageService()}
}

// POST /uploads/file — generic single-file upload
func (h *UploadHandler) UploadFile(c *gin.Context) {
	folder := c.DefaultPostForm("folder", "general")
	identifier := c.DefaultPostForm("identifier", "misc")

	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "file is required (multipart field 'file')", nil)
		return
	}

	if file.Size > maxFileSize {
		response.BadRequest(c, fmt.Sprintf("File too large (max %d MB)", maxFileSize>>20), nil)
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedImageExts[ext] && !allowedDocExts[ext] {
		response.BadRequest(c, "File type not allowed. Supported: images (jpg, png, gif, webp, svg), documents (pdf, doc, docx, xls, xlsx, ppt, pptx, txt, csv, zip, rar)", nil)
		return
	}

	url, size, mime, err := h.storage.UploadFile(file, folder, identifier)
	if err != nil {
		response.InternalServerError(c, "Upload failed", err.Error())
		return
	}

	fullURL := h.storage.ResolveURL(c.Request, url)

	response.Created(c, "File uploaded successfully", gin.H{
		"url":       fullURL,
		"file_name": file.Filename,
		"file_size": size,
		"mime_type": mime,
		"folder":    folder,
		"storage":   h.storageBackend(),
	})
}

// POST /uploads/image — image-only upload
func (h *UploadHandler) UploadImage(c *gin.Context) {
	folder := c.DefaultPostForm("folder", "images")
	identifier := c.DefaultPostForm("identifier", "misc")

	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "file is required (multipart field 'file')", nil)
		return
	}

	if file.Size > 10<<20 {
		response.BadRequest(c, "Image too large (max 10 MB)", nil)
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedImageExts[ext] {
		response.BadRequest(c, "Only image files allowed: jpg, jpeg, png, gif, webp, svg", nil)
		return
	}

	url, size, mime, err := h.storage.UploadFile(file, folder, identifier)
	if err != nil {
		response.InternalServerError(c, "Upload failed", err.Error())
		return
	}

	fullURL := h.storage.ResolveURL(c.Request, url)

	response.Created(c, "Image uploaded successfully", gin.H{
		"url":       fullURL,
		"file_name": file.Filename,
		"file_size": size,
		"mime_type": mime,
		"folder":    folder,
		"storage":   h.storageBackend(),
	})
}

// POST /uploads/document — document-only upload (pdf, office, etc.)
func (h *UploadHandler) UploadDocument(c *gin.Context) {
	folder := c.DefaultPostForm("folder", "documents")
	identifier := c.DefaultPostForm("identifier", "misc")

	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "file is required (multipart field 'file')", nil)
		return
	}

	if file.Size > maxFileSize {
		response.BadRequest(c, fmt.Sprintf("Document too large (max %d MB)", maxFileSize>>20), nil)
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedDocExts[ext] {
		response.BadRequest(c, "Only document files allowed: pdf, doc, docx, xls, xlsx, ppt, pptx, txt, csv, zip, rar", nil)
		return
	}

	url, size, mime, err := h.storage.UploadFile(file, folder, identifier)
	if err != nil {
		response.InternalServerError(c, "Upload failed", err.Error())
		return
	}

	fullURL := h.storage.ResolveURL(c.Request, url)

	response.Created(c, "Document uploaded successfully", gin.H{
		"url":       fullURL,
		"file_name": file.Filename,
		"file_size": size,
		"mime_type": mime,
		"folder":    folder,
		"storage":   h.storageBackend(),
	})
}

// DELETE /uploads — delete a file by URL
func (h *UploadHandler) DeleteFile(c *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "url is required in request body", nil)
		return
	}

	if err := h.storage.DeleteFile(req.URL); err != nil {
		response.InternalServerError(c, "Delete failed", err.Error())
		return
	}

	response.Success(c, "File deleted successfully", nil)
}

// GET /uploads/info — storage configuration info
func (h *UploadHandler) GetStorageInfo(c *gin.Context) {
	response.Success(c, "Storage configuration", gin.H{
		"backend":          h.storageBackend(),
		"max_file_size_mb": maxFileSize >> 20,
		"allowed_images":   []string{"jpg", "jpeg", "png", "gif", "webp", "svg"},
		"allowed_documents": []string{"pdf", "doc", "docx", "xls", "xlsx", "ppt", "pptx", "txt", "csv", "zip", "rar"},
	})
}

func (h *UploadHandler) storageBackend() string {
	if h.storage.IsS3() {
		return "s3"
	}
	return "local"
}
