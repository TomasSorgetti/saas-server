package di

import (
	"database/sql"
	"luthierSaas/internal/application/usecases/auth"
	"luthierSaas/internal/application/usecases/user"
	"luthierSaas/internal/infrastructure/cache"
	"luthierSaas/internal/infrastructure/config"
	"luthierSaas/internal/infrastructure/email"
	"luthierSaas/internal/infrastructure/logger"
	"luthierSaas/internal/infrastructure/persistance/repositories"
	"luthierSaas/internal/infrastructure/queue"
	"luthierSaas/internal/interfaces/http/handlers"

	"github.com/redis/go-redis/v9"
)

type Container struct {
	AuthHandler   *handlers.AuthHandler
	UserHandler   *handlers.UserHandler
	RedisClient   *redis.Client
	CacheService  *cache.Cache
}

func NewContainer(db *sql.DB, cfg *config.Config) (*Container, *email.EmailService) {
	// Logger
	log := logger.NewLogger()

	// Repositorios
	userRepo := repositories.NewUserRepository(db)
	suscriptionRepo := repositories.NewSubscriptionRepository(db)
	emailVerificationRepo := repositories.NewEmailVerificationRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)

	// Cliente Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	// Email Queue
	emailQueue := queue.NewQueue(redisClient, "email_queue")

	// Email Service
	emailService := email.NewEmailService(emailQueue)

	// Cache Service
	cacheService := cache.NewCache(redisClient)

	// Casos de uso
	authUC := auth.NewAuthUseCases(
		userRepo,
		suscriptionRepo,
		emailVerificationRepo,
		sessionRepo,
		emailService,
		cacheService,
		log,
		cfg.GoogleOAuth, // Inyectar la configuración de Google OAuth
	)
	userUC := user.NewUserUseCases(userRepo, sessionRepo, cacheService, emailService, log)

	// Handlers
	authHandler := handlers.NewAuthHandler(
		authUC.Login,
		authUC.Register,
		authUC.CheckEmail,
		authUC.VerifyEmail,
		authUC.ResendVerificationCode,
		authUC.RefreshToken,
		authUC.GoogleLogin,
		authUC.GoogleCallback,
		authUC.Logout,
	)
	userHandler := handlers.NewUserHandler(userUC.Profile, userUC.ChangePassword)

	return &Container{
		AuthHandler:  authHandler,
		UserHandler:  userHandler,
		RedisClient:  redisClient,
		CacheService: cacheService,
	}, emailService
}