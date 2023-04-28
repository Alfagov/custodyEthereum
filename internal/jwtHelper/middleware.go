package jwtHelper

import "github.com/gin-gonic/gin"

func JwtAuthMiddleware(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"message": "No Authorization Header"})
			return
		}

		_, role, err := ValidateAccessToken(authHeader)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"message": err.Error()})
			return
		}

		if !containsRole(roles, role) {
			c.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
			return
		}

		c.Next()
	}
}

func containsRole(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}
