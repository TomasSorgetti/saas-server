package middlewares

import (
	"luthierSaas/internal/infrastructure/security"
	"net/http"

	"github.com/gin-gonic/gin"
)

const RefreshUserIDKey = "userID"

func RefreshMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        cookie, err := c.Cookie("refresh_token") 
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid cookie"})
            return
        }

        userID, err := security.ValidateRefreshToken(cookie)

        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            return
        }

        c.Set(RefreshUserIDKey, userID)
        c.Next()
    }
}