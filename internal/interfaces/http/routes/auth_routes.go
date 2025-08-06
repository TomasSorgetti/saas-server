package routes

import (
	"log"
	"luthierSaas/internal/infrastructure/cache"
	"luthierSaas/internal/interfaces/http/handlers"
	"luthierSaas/internal/interfaces/http/middlewares"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
)

func SetupAuthRoutes(api *gin.RouterGroup, authHandler *handlers.AuthHandler, cacheService *cache.Cache) {
    signinLimiter, err := middlewares.NewRateLimiterMiddleware(cacheService, middlewares.RateLimiterConfig{
        Rate:   limiter.Rate{Period: time.Minute, Limit: 5},
        Prefix: "rate:auth:signin:free",
    })
    if err != nil {
        log.Fatalf("Failed to initialize signin rate limiter: %v", err)
    }

    signupLimiter, err := middlewares.NewRateLimiterMiddleware(cacheService, middlewares.RateLimiterConfig{
        Rate:   limiter.Rate{Period: time.Minute, Limit: 5},
        Prefix: "rate:auth:signup:free",
    })
    if err != nil {
        log.Fatalf("Failed to initialize signup rate limiter: %v", err)
    }

    checkEmailLimiter, err := middlewares.NewRateLimiterMiddleware(cacheService, middlewares.RateLimiterConfig{
        Rate:   limiter.Rate{Period: time.Minute, Limit: 10},
        Prefix: "rate:auth:check-email:free",
    })
    if err != nil {
        log.Fatalf("Failed to initialize check-email rate limiter: %v", err)
    }

    verifyEmailLimiter, err := middlewares.NewRateLimiterMiddleware(cacheService, middlewares.RateLimiterConfig{
        Rate:   limiter.Rate{Period: time.Minute, Limit: 10},
        Prefix: "rate:auth:verify-email:free",
    })
    if err != nil {
        log.Fatalf("Failed to initialize verify-email rate limiter: %v", err)
    }

    resendCodeLimiter, err := middlewares.NewRateLimiterMiddleware(cacheService, middlewares.RateLimiterConfig{
        Rate:   limiter.Rate{Period: time.Minute, Limit: 3},
        Prefix: "rate:auth:resend-code:free",
    })
    if err != nil {
        log.Fatalf("Failed to initialize resend-code rate limiter: %v", err)
    }

    refreshLimiter, err := middlewares.NewRateLimiterMiddleware(cacheService, middlewares.RateLimiterConfig{
        Rate:   limiter.Rate{Period: time.Minute, Limit: 10},
        Prefix: "rate:auth:refresh:free",
    })
    if err != nil {
        log.Fatalf("Failed to initialize refresh rate limiter: %v", err)
    }

    logoutLimiter, err := middlewares.NewRateLimiterMiddleware(cacheService, middlewares.RateLimiterConfig{
        Rate:   limiter.Rate{Period: time.Minute, Limit: 20},
        Prefix: "rate:auth:logout:free",
    })
    if err != nil {
        log.Fatalf("Failed to initialize logout rate limiter: %v", err)
    }

    auth := api.Group("/auth")
    {
        auth.POST("/signin", signinLimiter, authHandler.Login)
        auth.POST("/signup", signupLimiter, authHandler.Register)
        auth.POST("/check-email", checkEmailLimiter, authHandler.CheckEmail)
        auth.POST("/verify-email", verifyEmailLimiter, authHandler.VerifyEmail)
        auth.POST("/resend-code", resendCodeLimiter, authHandler.ResendVerificationCode)
        auth.POST("/refresh", middlewares.RefreshMiddleware(), refreshLimiter, authHandler.RefreshToken)
        auth.POST("/logout", logoutLimiter, authHandler.Logout)
        auth.GET("/google/login", authHandler.GoogleLogin)
        auth.GET("/google/callback", authHandler.GoogleCallback)
    }
}