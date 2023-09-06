// config/jwt.go
package config

import (
	"errors"
	"examn_go/infra/database"
	"examn_go/models"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateAccessToken(email string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = email
	claims["exp"] = time.Now().Add(time.Hour).Unix() // Access token expires in 1 hour
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func GenerateRefreshToken(email string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = email
	claims["exp"] = time.Now().Add(24 * time.Hour).Unix() // Refresh token expires in 24 hours
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
}

func TokenMatcher(accessToken string, email string) error {
	var userInfo models.User
	if err := database.DB.Select("access_token").Where("email = ?", email).First(&userInfo).Error; err != nil {
		return errors.New("failed to load student from DB")
	}

	if userInfo.AccessToken != accessToken {
		return errors.New("tokens do not match")
	}

	return nil
}
