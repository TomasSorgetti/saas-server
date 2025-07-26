package di

import (
	"database/sql"
	authUC "luthierSaas/internal/application/usecases/auth"
	"luthierSaas/internal/infrastructure/email"
	"luthierSaas/internal/infrastructure/persistance/repositories"
	"luthierSaas/internal/infrastructure/queue"
	"luthierSaas/internal/interfaces/http/handlers"

	"github.com/redis/go-redis/v9"
)

type Container struct {
	AuthHandler *handlers.AuthHandler
	UserHandler *handlers.UserHandler
}

func NewContainer(db *sql.DB) (*Container, *email.EmailService) {
	// Repositorios
	userRepo := repositories.NewUserRepository(db)
	emailVerificationRepo := repositories.NewEmailVerificationRepository(db)

	// Client Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	// Email Queue
	emailQueue := queue.NewQueue(redisClient, "email_queue")

	// Email Service
	emailService := email.NewEmailService(emailQueue)

	// Use cases
	authUC := authUC.NewAuthUseCases(userRepo, emailVerificationRepo, emailService)

	// Handlers
	authHandler := handlers.NewAuthHandler(authUC.Register, authUC.CheckEmail, authUC.VerifyEmail)

	return &Container{
		AuthHandler: authHandler,
	}, emailService
}
