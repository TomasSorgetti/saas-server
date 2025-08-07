package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"luthierSaas/internal/infrastructure/cache"
	"time"

	"golang.org/x/oauth2"
)

type LoginGoogleUseCase struct {
	oauthConfig  *oauth2.Config
	cacheService *cache.Cache
}

func NewLoginGoogleUseCase(oauthConfig *oauth2.Config, cacheService *cache.Cache) *LoginGoogleUseCase {
	return &LoginGoogleUseCase{
		oauthConfig:  oauthConfig,
		cacheService: cacheService,
	}
}

func (uc *LoginGoogleUseCase) Execute(ctx context.Context) (string, error) {
	state, err := generateRandomState()
	if err != nil {
		return "", fmt.Errorf("failed to generate state: %w", err)
	}

	err = uc.cacheService.Set(ctx, "oauth_state:"+state, state, 10*time.Minute)
	if err != nil {
		return "", fmt.Errorf("failed to store state in cache: %w", err)
	}

	url := uc.oauthConfig.AuthCodeURL(state)
	return url, nil
}

func generateRandomState() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}