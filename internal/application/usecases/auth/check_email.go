package auth

import (
	"context"
	"fmt"
	"luthierSaas/internal/infrastructure/cache"
	"luthierSaas/internal/interfaces/repository"
	"time"
)

type CheckEmailUseCase struct {
	userRepo repository.UserRepository
	cacheService *cache.Cache
}

func NewCheckEmailUseCase(userRepo repository.UserRepository, cacheService *cache.Cache) *CheckEmailUseCase {
	return &CheckEmailUseCase{
		userRepo,
		cacheService,
	}
}

func (uc *CheckEmailUseCase) Execute(ctx context.Context, email string) (bool, error) {
    cacheKey := fmt.Sprintf("email:check:%s", email)

    cached, err := uc.cacheService.Get(ctx, cacheKey)
    if err == nil && cached != "" {
        switch cached {
            case "true":
                return true, nil
            case "false":
                return false, nil
        }
    }

    exists, err := uc.userRepo.EmailExists(email)
    if err != nil {
        return false, err
    }

    cacheValue := "false"
    if exists {
        cacheValue = "true"
    }

    err = uc.cacheService.Set(ctx, cacheKey, cacheValue, 5*time.Minute)
    if err != nil {
        fmt.Printf("Failed to set cache for %s: %v\n", cacheKey, err)
    }

    return exists, nil
}