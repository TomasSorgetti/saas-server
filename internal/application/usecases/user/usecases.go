package user

import (
	"luthierSaas/internal/infrastructure/cache"
	"luthierSaas/internal/interfaces/repository"
)

type UserUseCases struct {
    Profile      *ProfileUseCase
}

func NewUserUseCases(userRepo repository.UserRepository, cacheService *cache.Cache) *UserUseCases{
    return &UserUseCases{
        Profile:      NewProfileUseCase(userRepo, cacheService ),
    }
}