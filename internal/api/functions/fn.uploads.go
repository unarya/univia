package functions

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/unarya/univia/internal/api/modules/post/models"
	"github.com/unarya/univia/pkg/utils"
)

// Allowed Media Types
var allowedMediaTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
	"video/mp4":  true,
	"video/avi":  true,
	"video/mov":  true,
}

func SaveMediaToServer(files []*multipart.FileHeader) ([]posts.Media, *utils.ServiceError) {
	// Handle file uploads
	var savedMedia []posts.Media

	for _, file := range files {
		fileType := file.Header.Get("Content-Type")

		// Validate file type
		if !allowedMediaTypes[fileType] {
			return nil, &utils.ServiceError{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid file format. Allowed: JPEG, PNG, GIF, MP4, AVI, MOV",
			}
		}

		// Generate unique filename
		uniqueFileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
		mediaPath := filepath.Join("uploads", uniqueFileName)

		// Save file
		if err := SaveUploadedFile(file, mediaPath); err != nil {
			return nil, &utils.ServiceError{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to save file",
			}
		}

		// Append to saved media list
		savedMedia = append(savedMedia, posts.Media{Path: mediaPath, Type: fileType})
	}

	return savedMedia, nil
}

// SaveUploadedFile saves the uploaded file to disk
func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

// DeleteUploadedFile deletes a file from disk
func DeleteUploadedFile(dst string) error {
	if err := os.Remove(dst); err != nil {
		return err
	}
	return nil
}
