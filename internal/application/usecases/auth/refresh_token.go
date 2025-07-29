package auth

import (
	"luthierSaas/internal/infrastructure/security"
	"luthierSaas/internal/interfaces/http/dtos"
)

type RefreshTokenUseCase struct {
}

func NewRefreshTokenUseCase() *RefreshTokenUseCase {
	return &RefreshTokenUseCase{}
}

func (uc *RefreshTokenUseCase) Execute(userID int) (*dtos.RefreshResponse, error) {
	accessToken, err := security.CreateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := security.CreateRefreshToken(userID)
	if err != nil {
		return nil, err
	}

	return &dtos.RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}