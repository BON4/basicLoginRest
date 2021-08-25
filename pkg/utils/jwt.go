package utils

import (
	"basicLoginRest/config"
	"basicLoginRest/internal/models"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Claims struct {
	ID    	 uint `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

func GenerateJWTToken(cfg *config.Config, user *models.User) (string, error) {
	claims := &Claims{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second*60).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(cfg.Server.JwtSecretKey))
	if err != nil {
		return "", err
	}
	return token, err
}
