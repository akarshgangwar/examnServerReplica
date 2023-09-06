package middleware

import (
	"examn_go/config"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		tokenString := GetTokenFromHeader(authHeader)
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{"error":"token is missing"})
			return
		}
		token, err := config.ValidateToken(tokenString) // Validate the token
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is invalid"})
			return
		}
		claims, _ := token.Claims.(jwt.MapClaims)
		userEmail := claims["sub"].(string)
		c.Set("email",userEmail)

		err = config.TokenMatcher(tokenString, userEmail) // Check if the access token matches the one stored in the database
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "access token doesn't matched"})
			return
		}
		c.Next()
	}
}

func GetTokenFromHeader(authHeader string) string {
	const bearerPrefix = "Bearer "
	if strings.HasPrefix(authHeader, bearerPrefix) {
		return authHeader[len(bearerPrefix):]
	}
	return ""
}
