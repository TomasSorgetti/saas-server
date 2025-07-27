package repository

import (
	"context"
	"luthierSaas/internal/domain/entities"
)


type EmailVerificationRepository interface {
	Create(ctx context.Context, ev *entities.EmailVerification) error
	GetByUserID(ctx context.Context, userID int) (*entities.EmailVerification, error)
	MarkAsVerified(ctx context.Context, userID int) error
	DeleteByUserID(ctx context.Context, userID int) error
	UpdateCode(ctx context.Context, id int, newCode string) error
}