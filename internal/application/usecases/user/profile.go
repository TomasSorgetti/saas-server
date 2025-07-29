package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"luthierSaas/internal/infrastructure/cache"
	"luthierSaas/internal/interfaces/http/dtos"
	"luthierSaas/internal/interfaces/repository"
	"time"
)

type ProfileUseCase struct {
	userRepo repository.UserRepository
	cache    *cache.Cache
}

func NewProfileUseCase(userRepo repository.UserRepository, cache *cache.Cache) *ProfileUseCase {
	return &ProfileUseCase{userRepo, cache}
}

func (uc *ProfileUseCase) Execute(userID int) (*dtos.ProfileResponse, error) {
	cacheKey := fmt.Sprintf("profile:user:%d", userID)
	ctx := context.Background()

	cachedProfile, err := uc.cache.Get(ctx, cacheKey)
	if err == nil {
		var profile dtos.ProfileResponse
		if err := json.Unmarshal([]byte(cachedProfile), &profile); err == nil {
			return &profile, nil
		}
		fmt.Printf("Error deserializing cached profile for user %d: %v\n", userID, err)
	}

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
		Subscription: user.Subscription,
	}

	profileJSON, err := json.Marshal(profile)
	if err != nil {
		fmt.Printf("Error serializing profile for user %d: %v\n", userID, err)
	} else {
		ttl := 5 * time.Minute
		if err := uc.cache.Set(ctx, cacheKey, string(profileJSON), ttl); err != nil {
			fmt.Printf("Error caching profile for user %d: %v\n", userID, err)
		}
	}

	return profile, nil
}