package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"gin-demo/pkg/storage"

	"github.com/gin-gonic/gin"
)

// FileUploadHandler handles file upload operations
type FileUploadHandler struct {
	storage *storage.S3Client
}

// NewFileUploadHandler creates a new file upload handler
func NewFileUploadHandler(storage *storage.S3Client) *FileUploadHandler {
	return &FileUploadHandler{storage: storage}
}

// UploadFile handles file upload
func (h *FileUploadHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"code":    "INVALID_FILE",
			"message": "No file provided",
		})
		return
	}

	// Validate file size (10MB limit)
	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"code":    "FILE_TOO_LARGE",
			"message": "File size exceeds 10MB limit",
		})
		return
	}

	// Generate unique key
	timestamp := time.Now().Unix()
	ext := filepath.Ext(file.Filename)
	key := fmt.Sprintf("uploads/%d_%s%s", timestamp, file.Filename, ext)

	// Open file
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"code":    "INTERNAL_ERROR",
			"message": "Failed to open file",
		})
		return
	}
	defer src.Close()

	// Upload to S3
	if err := h.storage.Upload(c.Request.Context(), key, src, file.Header.Get("Content-Type")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"code":    "UPLOAD_FAILED",
			"message": "Failed to upload file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "File uploaded successfully",
		"data": gin.H{
			"key":      key,
			"filename": file.Filename,
			"size":     file.Size,
		},
	})
}

// GetPresignedURL generates a presigned URL for file download
func (h *FileUploadHandler) GetPresignedURL(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"code":    "MISSING_KEY",
			"message": "File key is required",
		})
		return
	}

	// Generate presigned URL (valid for 1 hour)
	url, err := h.storage.GetPresignedURL(c.Request.Context(), key, 1*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"code":    "URL_GENERATION_FAILED",
			"message": "Failed to generate download URL",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"url":        url,
			"expires_in": "1h",
		},
	})
}

// ListFiles lists uploaded files
func (h *FileUploadHandler) ListFiles(c *gin.Context) {
	prefix := c.DefaultQuery("prefix", "uploads/")

	files, err := h.storage.List(c.Request.Context(), prefix)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"code":    "LIST_FAILED",
			"message": "Failed to list files",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"files": files,
			"count": len(files),
		},
	})
}

// DeleteFile deletes a file
func (h *FileUploadHandler) DeleteFile(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"code":    "MISSING_KEY",
			"message": "File key is required",
		})
		return
	}

	if err := h.storage.Delete(c.Request.Context(), key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"code":    "DELETE_FAILED",
			"message": "Failed to delete file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "File deleted successfully",
	})
}
