package user

import (
	"luthierSaas/internal/infrastructure/cache"
	"luthierSaas/internal/infrastructure/email"
	"luthierSaas/internal/interfaces/repository"

	"github.com/rs/zerolog"
)

type UserUseCases struct {
    Profile      *ProfileUseCase
    ChangePassword *ChangePasswordUseCase
}

func NewUserUseCases(
    userRepo repository.UserRepository, 
    sessionRepo repository.SessionRepository, 
    cacheService *cache.Cache, 
    emailService *email.EmailService,
    logger      *zerolog.Logger) *UserUseCases{
        
    return &UserUseCases{
        Profile:      NewProfileUseCase(userRepo, sessionRepo, cacheService ),
        ChangePassword: NewChangePasswordUseCase(userRepo, sessionRepo, cacheService, emailService),
    }
}