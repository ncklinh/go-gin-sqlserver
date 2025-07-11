package middleware

import (
	token "film-rental/internal/token"
	"film-rental/internal/token/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequirePermission creates a middleware that checks if the user has a specific permission
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the payload from the auth middleware
		payloadInterface, exists := c.Get(authorizationPayloadKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization payload not found"})
			return
		}

		payload, ok := payloadInterface.(*token.Payload)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid authorization payload"})
			return
		}

		// Check if the user has the required permission
		if !model.HasPermission(payload.Role, permission) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":               "Insufficient permissions",
				"message":             "You don't have permission to perform this action",
				"required_permission": permission,
				"user_role":           payload.Role,
			})
			return
		}

		c.Next()
	}
}
