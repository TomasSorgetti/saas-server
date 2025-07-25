package routes

import (
	"database/sql"
	"luthierSaas/internal/di"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, db *sql.DB, container *di.Container) {
    api := r.Group("/api/v1")

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
    // auth routes
    SetupAuthRoutes(api, container.AuthHandler)
}