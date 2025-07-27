package middlewares

import (
	"luthierSaas/internal/infrastructure/security"
	"net/http"

	"github.com/gin-gonic/gin"
)

const UserIDKey = "userID"

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        cookie, err := c.Cookie("access_token") 
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid cookie"})
            return
        }

        // Validar el token almacenado en la cookie
        userID, err := security.ValidateAccessToken(cookie)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            return
        }

        c.Set(UserIDKey, userID)
        c.Next()
    }
}