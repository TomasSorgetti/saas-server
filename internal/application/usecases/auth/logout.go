package auth

import (
	"context"
	"fmt"
	"luthierSaas/internal/infrastructure/security"
	"luthierSaas/internal/interfaces/repository"
)

type LogoutUseCase struct {
	sessionRepo  repository.SessionRepository
}


func NewLogoutUseCase(sessionRepo repository.SessionRepository) *LogoutUseCase {
	return &LogoutUseCase{
		sessionRepo: sessionRepo,
	}
}

func (uc *LogoutUseCase) Execute(ctx context.Context, accessToken string) ( error) {
    accessTokenHash, err := security.HashToken(accessToken)
    if err != nil {
        return fmt.Errorf("failed to hash access token: %w", err)
    }
	
    err = uc.sessionRepo.Delete(ctx, accessTokenHash)
    if err != nil {
        return fmt.Errorf("failed to delete session: %w", err)
    }

    return nil
}