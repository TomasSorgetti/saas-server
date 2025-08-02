package user

import (
	"luthierSaas/internal/infrastructure/cache"
	"luthierSaas/internal/infrastructure/email"
	"luthierSaas/internal/interfaces/repository"
)

type UserUseCases struct {
    Profile      *ProfileUseCase
    ChangePassword *ChangePasswordUseCase
}

func NewUserUseCases(userRepo repository.UserRepository, cacheService *cache.Cache, emailService *email.EmailService) *UserUseCases{
    return &UserUseCases{
        Profile:      NewProfileUseCase(userRepo, cacheService ),
        ChangePassword: NewChangePasswordUseCase(userRepo, cacheService, emailService),
    }
}