package routes

import (
	"luthierSaas/internal/interfaces/http/handlers"
	"luthierSaas/internal/interfaces/http/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(api *gin.RouterGroup, userHandler *handlers.UserHandler) {

    users := api.Group("/users")
    {
        users.GET("profile", middlewares.AuthMiddleware(), userHandler.GetProfile)
    }
}