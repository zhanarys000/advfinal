package middleware

import (
	"adv/pkg/logger"
	"adv/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func RequireAuthMiddleware(c *gin.Context) {
	logger := logger.GetLogger()

	var token string
	authHeader := c.GetHeader("Authorization")

	fields := strings.Fields(authHeader)
	if len(fields) != 0 && fields[0] == "Bearer" {
		token = fields[1]
	}

	if authHeader == "" {
		logger.Error("Authorization header not found")
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Try to sign in first"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	fmt.Println(token)

	if token == "" {
		logger.Error("Invalid token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	id, email, userType, err := utils.VerifyToken(token)
	if err != nil {
		logger.Error("Token verification failed:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token verification failed"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("id", id)
	c.Set("email", email)
	c.Set("userType", userType)

	logger.Info("Authentication successful")
	c.Next()
}
func AdminMiddleware(c *gin.Context) {
	userType, ok := c.Get("userType")
	if !ok {
		logger := logger.GetLogger()
		logger.Error("User type not found in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User type not found"})
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	userTypeStr, ok := userType.(string)
	if !ok {
		logger := logger.GetLogger()
		logger.Error("Invalid user type format")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type format"})
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if userTypeStr != "ADMIN" {
		logger := logger.GetLogger()
		logger.Error("User is not an admin")
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this resource"})
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	logger := logger.GetLogger()
	logger.Info("Admin authorization successful")
	c.Next()
}
