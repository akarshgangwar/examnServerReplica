package controllers

import (
	"examn_go/config"
	"examn_go/models"
	"examn_go/repository"
	"net/http"
	"strings"

	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	LoggedInFlag bool   `json:"anywayLogin"`
}

func LoginHandler(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if req.Email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Username or Email is required"})
		return
	}
	student, err := repository.GetUserByEmail(req.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	if req.Password != student.Password {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials pass"})
		return
	}
	if !req.LoggedInFlag {
		currentUTCTime := time.Now().UTC()
		expirationTimeUTC := student.TokenExpirationTime.UTC()
		if student.AccessToken != "" && expirationTimeUTC.After(currentUTCTime) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "User already logged in"})
			return
		}
	}
	// Generate JWT access token
	accessToken, err := config.GenerateAccessToken(req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}
	// Generate JWT refresh token
	refreshToken, err := config.GenerateRefreshToken(req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}
	student.AccessToken = accessToken
	student.RefreshToken = refreshToken
	student.TokenExpirationTime = time.Now().Add(time.Hour)

	err = repository.Update(&student)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save tokens"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"access Token: ": accessToken, "refresh Token: ": refreshToken})
}

func RefreshTokensHandler(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error ": " Refresh Token is missing"})
	}
	refreshToken := authHeader[len(bearerPrefix):]
	token, err := config.ValidateToken(refreshToken) // Validate the token
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if !token.Valid {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is invalid"})
		return
	}
	// extracting email from claims
	claims, _ := token.Claims.(jwt.MapClaims)
	userEmail := claims["sub"].(string)
	// Generate JWT access token
	accessToken, err := config.GenerateAccessToken(userEmail)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}
	// Generate JWT refresh token
	refreshToken, err = config.GenerateRefreshToken(userEmail)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}
	userTokenInfo := models.User{
		AccessToken:         accessToken,
		RefreshToken:        refreshToken,
		TokenExpirationTime: time.Now().Add(time.Hour),
		Email:               userEmail,
	}
	if err := repository.Update(&userTokenInfo); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save tokens"})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{"access Token: ": accessToken, "refresh Token: ": refreshToken})
}
