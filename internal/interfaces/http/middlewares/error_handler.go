package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"luthierSaas/internal/interfaces/http/errors"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			if httpErr, ok := err.(*errors.HTTPError); ok {
				c.JSON(httpErr.Code, gin.H{
					"success": false,
					"error":   httpErr.Message,
					"details": httpErr.Details,
				})
				return
			}

			// Fallback: Error no manejado
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Ocurrió un error inesperado",
				"details": err.Error(),
			})
		}
	}
}