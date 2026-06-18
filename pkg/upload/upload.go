package upload

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
}

const maxFileSize = 2 * 1024 * 1024 // 2MB

func SavePhoto(file *multipart.FileHeader, storagePath string) (string, error) {
	if file.Size > maxFileSize {
		return "", fmt.Errorf("file size exceeds 2MB limit")
	}

	contentType := file.Header.Get("Content-Type")
	if !allowedImageTypes[contentType] {
		return "", fmt.Errorf("file type not allowed, use jpeg, png, or webp")
	}

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)

	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return "", fmt.Errorf("failed to create storage directory: %w", err)
	}

	dst := filepath.Join(storagePath, filename)
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer out.Close()

	buf := make([]byte, 32*1024)
	for {
		n, err := src.Read(buf)
		if n > 0 {
			if _, werr := out.Write(buf[:n]); werr != nil {
				return "", werr
			}
		}
		if err != nil {
			break
		}
	}

	return filename, nil
}

func DeletePhoto(storagePath, filename string) error {
	if filename == "" {
		return nil
	}
	path := filepath.Join(storagePath, filepath.Base(filename))
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func BuildPhotoURL(storageURL, filename string) string {
	if filename == "" {
		return ""
	}
	return strings.TrimRight(storageURL, "/") + "/" + filename
}
