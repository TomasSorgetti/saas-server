package auth

import (
	"context"
	"fmt"
	"luthierSaas/internal/infrastructure/security"
	"luthierSaas/internal/interfaces/repository"

	"github.com/rs/zerolog"
)

type LogoutUseCase struct {
	sessionRepo  repository.SessionRepository
    logger      *zerolog.Logger
}


func NewLogoutUseCase(
    sessionRepo repository.SessionRepository,
    logger *zerolog.Logger) *LogoutUseCase {

	return &LogoutUseCase{
		sessionRepo: sessionRepo,
        logger: logger,
	}
}

func (uc *LogoutUseCase) Execute(ctx context.Context, accessToken string) ( error) {
    accessTokenHash, err := security.HashToken(accessToken)
    if err != nil {
        uc.logger.Error().
            Err(err).
            Str("access_token_hash", accessTokenHash).
            Msg("Failed to hash access token")
        return fmt.Errorf("failed to hash access token: %w", err)
    }
	
    err = uc.sessionRepo.Delete(ctx, accessTokenHash)
    if err != nil {
        uc.logger.Error().
            Err(err).
            Str("access_token_hash", accessTokenHash).
            Msg("Failed to delete session")
        return fmt.Errorf("failed to delete session: %w", err)
    }

    uc.logger.Info().
        Str("access_token_hash", accessTokenHash).
        Msg("Session deleted successfully")
    return nil
}