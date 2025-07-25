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

func NewContainer(db *sql.DB) *Container {
	// Repositorios
	userRepo := repositories.NewMySQLUserRepository(db)

	// Cliente Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: "redis:6379", // Cambiar si es otro host/puerto
	})

	// Cola para emails
	emailQueue := queue.NewQueue(redisClient, "email_queue")

	// Servicio de Email
	emailService := email.NewEmailService(emailQueue)

	// Use cases
	authUC := authUC.NewAuthUseCases(userRepo, emailService)

	// Handlers
	authHandler := handlers.NewAuthHandler(authUC.Register, authUC.CheckEmail)

	return &Container{
		AuthHandler: authHandler,
	}
}
