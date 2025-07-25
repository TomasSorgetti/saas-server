package routes

import (
	"luthierSaas/internal/interfaces/http/handlers"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(api *gin.RouterGroup) {
    userHandler := handlers.NewUserHandler()

    users := api.Group("/users")
    {
        users.POST("", userHandler.Create)
        users.GET("", userHandler.GetAll)
        users.GET("/:id", userHandler.GetByID)
        users.PUT("/:id", userHandler.Update)
        users.DELETE("/:id", userHandler.Delete)
    }
}