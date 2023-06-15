package service

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type AvatarService struct{}

func NewAvatarService() *AvatarService {
	return &AvatarService{}
}

func (s *AvatarService) GenerateAvatar(userID int) (string, error) {

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using environment variables")
	}

	avatarDir := "avatars"
	if _, err := os.Stat(avatarDir); os.IsNotExist(err) {
		err := os.MkdirAll(avatarDir, 0755)
		if err != nil {
			return "", err
		}
	}

	avatarPath := filepath.Join(avatarDir, fmt.Sprintf("%d.png", userID))

	resp, err := http.Get(os.Getenv("AVATAR_GEN_URL"))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	out, err := os.Create(avatarPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return avatarPath, nil
}
