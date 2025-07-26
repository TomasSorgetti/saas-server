package repository

import (
	"context"
	"luthierSaas/internal/domain/entities"
)


type EmailVerificationRepository interface {
	Create(ctx context.Context, ev *entities.EmailVerification) error
	GetByUserID(ctx context.Context, userID int64) (*entities.EmailVerification, error)
	MarkAsVerified(ctx context.Context, userID int64) error
	DeleteByUserID(ctx context.Context, userID int64) error
}