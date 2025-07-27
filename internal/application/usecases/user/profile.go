package user

import (
	"luthierSaas/internal/infrastructure/cache"
	"luthierSaas/internal/interfaces/repository"
)

type ProfileUseCase struct {
	userRepo repository.UserRepository
	cache    *cache.Cache
}

func NewProfileUseCase(userRepo repository.UserRepository, cache *cache.Cache) *ProfileUseCase {
	return &ProfileUseCase{userRepo, cache}
}

func (uc *ProfileUseCase) Execute(userID int) (bool, error) {
	
	user, err := uc.userRepo.FindByID(userID)
	
	if err != nil {
		return false, err
	}
	
	return user != nil, nil
}