package utils

import (
	"mime/multipart"
	"net/http"
)

const (
	maxFileSize = 1024 * 1024 * 10
)

func CheckImageUpload(file multipart.File) (string, error) {
	fileHeader := make([]byte, maxFileSize)
	ContentType := ""
	if _, err := file.Read(fileHeader); err != nil {
		return ContentType, err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return ContentType, err
	}

	count, err := file.Seek(0, 2)
	if err != nil {
		return ContentType, err
	}
	if count > maxFileSize {
		return ContentType, err
	}
	if _, err := file.Seek(0, 0); err != nil {
		return ContentType, err
	}
	ContentType = http.DetectContentType(fileHeader)

	if ContentType != "image/jpg" && ContentType != "image/png" && ContentType != "image/jpeg" {
		return ContentType, err
	}

	return ContentType, nil
}
