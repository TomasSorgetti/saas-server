package auth

import (
	"context"
	"errors"
	"fmt"
	"luthierSaas/internal/domain/entities"
	"luthierSaas/internal/infrastructure/security"
	"luthierSaas/internal/interfaces/http/dtos"
	"luthierSaas/internal/interfaces/repository"
	"time"
)

type RefreshTokenUseCase struct {
	userRepo repository.UserRepository
	sessionRepo repository.SessionRepository
}

func NewRefreshTokenUseCase(userRepo repository.UserRepository, sessionRepo repository.SessionRepository) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{userRepo, sessionRepo}
}

func (uc *RefreshTokenUseCase) Execute(ctx context.Context, refreshToken string, deviceInfo string) (*dtos.RefreshResponse, error) {
    refreshTokenHash, err := security.HashToken(refreshToken)
    if err != nil {
        return nil, fmt.Errorf("failed to hash refresh token: %w", err)
    }

    session, err := uc.sessionRepo.FindByRefreshTokenHash(ctx, refreshTokenHash)
    if err != nil {
        return nil, fmt.Errorf("failed to find session: %w", err)
    }
    if session == nil {
        return nil, errors.New("invalid or expired refresh token")
    }

    user, err := uc.userRepo.FindByID(session.UserID)
    if err != nil {
        return nil, fmt.Errorf("failed to find user: %w", err)
    }
    if user == nil {
        return nil, errors.New("user not found")
    }
    if user.Deleted {
        return nil, errors.New("user deleted")
    }
    if !user.Verified {
        return nil, errors.New("user not verified")
    }

    accessToken, err := security.CreateAccessToken(user.ID)
    if err != nil {
        return nil, fmt.Errorf("failed to create access token: %w", err)
    }
    newRefreshToken, err := security.CreateRefreshToken(user.ID)
    if err != nil {
        return nil, fmt.Errorf("failed to create refresh token: %w", err)
    }

    accessTokenHash, err := security.HashToken(accessToken)
    if err != nil {
        return nil, fmt.Errorf("failed to hash access token: %w", err)
    }
    newRefreshTokenHash, err := security.HashToken(newRefreshToken)
    if err != nil {
        return nil, fmt.Errorf("failed to hash refresh token: %w", err)
    }

    newSession := &entities.Session{
        UserID:           user.ID,
        AccessTokenHash:  accessTokenHash,
        RefreshTokenHash: newRefreshTokenHash,
        ExpiresAt:        time.Now().Add(15 * time.Minute),
        RefreshExpiresAt: time.Now().Add(7 * 24 * time.Hour),
        IsValid:          true,
        DeviceInfo:       deviceInfo, 
    }

    err = uc.sessionRepo.Create(ctx, newSession)
    if err != nil {
        return nil, fmt.Errorf("failed to create new session: %w", err)
    }

    err = uc.sessionRepo.Delete(ctx, session.AccessTokenHash)
    if err != nil {
        return nil, fmt.Errorf("failed to delete old session: %w", err)
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

    return &dtos.RefreshResponse{
        Profile:      profile,
        AccessToken:  accessToken,
        RefreshToken: newRefreshToken,
    }, nil
}