package routes

import (
	"luthierSaas/internal/interfaces/http/handlers"

	"github.com/gin-gonic/gin"
)
func SetupAuthRoutes(api *gin.RouterGroup, authHandler *handlers.AuthHandler) {
    auth := api.Group("/auth")
    {
        auth.POST("/signin", authHandler.Login)
        auth.POST("/signup", authHandler.Register)
        auth.POST("/check-email", authHandler.CheckEmail)
        auth.POST("/verify-email", authHandler.VerifyEmail)
        auth.POST("/resend-code", authHandler.ResendVerificationCode)
    }
}