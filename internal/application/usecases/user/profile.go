package user

import (
	"errors"
	"luthierSaas/internal/infrastructure/cache"
	"luthierSaas/internal/interfaces/http/dtos"
	"luthierSaas/internal/interfaces/repository"
)

type ProfileUseCase struct {
	userRepo repository.UserRepository
	cache    *cache.Cache
}

func NewProfileUseCase(userRepo repository.UserRepository, cache *cache.Cache) *ProfileUseCase {
	return &ProfileUseCase{userRepo, cache}
}

func (uc *ProfileUseCase) Execute(userID int) (*dtos.ProfileResponse, error) {
	
	user, err := uc.userRepo.FindByID(userID)
    if err != nil {
        return nil, errors.New("user not found")
    }

	if user.Deleted {
        return nil, errors.New("user deleted")
    }

   if !user.Verified {
    return nil, errors.New("user not verified")
}
	
	profile := &dtos.ProfileResponse{
        ID:           user.ID,
        Email:        user.Email,
        FirstName:    user.FirstName,
        LastName:     user.LastName,
        Phone:        user.Phone,
        Country:      user.Country,
        WorkshopName: user.WorkshopName,
        LastLogin:    user.LastLogin,
    }

    return profile, nil
}