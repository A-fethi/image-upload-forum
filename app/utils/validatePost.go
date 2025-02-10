package utils

import (
	"fmt"
	"forum/app/config"
	"path/filepath"
)

func ValidatePost(title, content string) error {
	if title == "" {
		config.Logger.Println("Validation failed: Title is empty")
		return fmt.Errorf("title cannot be empty")
	}

	if content == "" {
		config.Logger.Println("Validation failed: Content is empty")
		return fmt.Errorf("content cannot be empty")
	}

	if len(title) > 100 {
		config.Logger.Printf("Validation failed: Title length is invalid (title length: %d)", len(title))
		return fmt.Errorf("title must be between 5 and 100 characters")
	}

	if len(content) > 5000 {
		config.Logger.Printf("Validation failed: Content length is invalid (content length: %d)", len(content))
		return fmt.Errorf("content must be between 10 and 5000 characters")
	}

	config.Logger.Println("Validation successful for title and content")
	return nil
}

func ValidateImage(imageName string) error {
	ext := filepath.Ext(imageName)
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".svg":  true,
	}

	if !allowedExtensions[ext] {
		config.Logger.Println("Validation failed: wrong image extension", ext)
		return fmt.Errorf("image extension is not allowed")
	}
	config.Logger.Println("Validation successful for image:", imageName)
	return nil
}
