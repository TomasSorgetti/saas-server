package repository

import (
	"context"
	"luthierSaas/internal/domain/entities"
)

type SessionRepository interface {
    Create(ctx context.Context, session *entities.Session) error
    FindByAccessTokenHash(ctx context.Context, accessTokenHash string) (*entities.Session, error)
    FindByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*entities.Session, error)
    Update(ctx context.Context, session *entities.Session) error
    Invalidate(ctx context.Context, accessTokenHash string) error
    FindByUserID(ctx context.Context, userID int64) ([]*entities.Session, error)
}