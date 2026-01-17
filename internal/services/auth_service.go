package services

import (
	"errors"
	"time"

	"rcs-onboarding/internal/config"
	"rcs-onboarding/internal/repositories"
	"rcs-onboarding/internal/utils"

	"github.com/golang-jwt/jwt/v4"
)

type AuthService struct {
	repo *repositories.UserRepo
}

func NewAuthService(repo *repositories.UserRepo) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Login(username, password string) (string, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	cfg := config.LoadConfig()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString(cfg.JWTKey)
}
