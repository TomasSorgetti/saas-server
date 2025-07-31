package di

import (
	"database/sql"
	authUseCases "luthierSaas/internal/application/usecases/auth"
	userUseCases "luthierSaas/internal/application/usecases/user"
	"luthierSaas/internal/infrastructure/cache"
	"luthierSaas/internal/infrastructure/email"
	"luthierSaas/internal/infrastructure/persistance/repositories"
	"luthierSaas/internal/infrastructure/queue"
	"luthierSaas/internal/interfaces/http/handlers"

	"github.com/redis/go-redis/v9"
)

type Container struct {
	AuthHandler *handlers.AuthHandler
	UserHandler *handlers.UserHandler
	RedisClient *redis.Client
	CacheService *cache.Cache
}

func NewContainer(db *sql.DB) (*Container, *email.EmailService) {
	// Repositorios
	userRepo := repositories.NewUserRepository(db)
	suscriptionRepo := repositories.NewSubscriptionRepository(db)
	emailVerificationRepo := repositories.NewEmailVerificationRepository(db)

	// Client Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	// Email Queue
	emailQueue := queue.NewQueue(redisClient, "email_queue")

	// Email Service
	emailService := email.NewEmailService(emailQueue)
	
	// Cache Service
	cacheService := cache.NewCache(redisClient)

	// Use cases
	authUC := authUseCases.NewAuthUseCases(userRepo, suscriptionRepo, emailVerificationRepo, emailService, cacheService)
	userUC := userUseCases.NewUserUseCases(userRepo, cacheService)

	// Handlers
	authHandler := handlers.NewAuthHandler( authUC.Login, authUC.Register, authUC.CheckEmail, authUC.VerifyEmail, authUC.ResendVerificationCode, authUC.RefreshToken)
	
	userHandler := handlers.NewUserHandler(userUC.Profile)

	return &Container{
		AuthHandler: authHandler,
		UserHandler: userHandler,
		RedisClient: redisClient,
		CacheService: cacheService,
	}, emailService
}
