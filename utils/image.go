package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const ProductsDir = "media/products"

func SaveProductImage(fileHeader *multipart.FileHeader) (string, error) {
	if err := os.MkdirAll(ProductsDir, os.ModePerm); err != nil {
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowed[ext] {
		return "", fmt.Errorf("extensión de archivo no permitida: %s", ext)
	}

	filename := uuid.New().String() + ext
	dstPath := filepath.Join(ProductsDir, filename)

	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return "products/" + filename, nil
}

func DeleteProductImage(relativePath string) error {
	if relativePath == "" {
		return nil
	}
	fullPath := filepath.Join("media", relativePath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil // ya no existe, no es error
	}
	return os.Remove(fullPath)
}
