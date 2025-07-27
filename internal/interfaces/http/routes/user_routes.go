package routes

import (
	"luthierSaas/internal/interfaces/http/handlers"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(api *gin.RouterGroup, userHandler *handlers.UserHandler) {

    users := api.Group("/users")
    {
        users.GET("profile", userHandler.GetProfile)
    }
}