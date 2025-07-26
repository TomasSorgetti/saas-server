package auth

import (
	"context"
	"errors"
	"luthierSaas/internal/infrastructure/security"
	"luthierSaas/internal/interfaces/repository"
	"time"
)

type VerifyEmailUseCase struct {
	userRepo              repository.UserRepository
	emailVerificationRepo  repository.EmailVerificationRepository
}

func NewVerifyEmailUseCase(userRepo repository.UserRepository, emailVerificationRepo repository.EmailVerificationRepository) *VerifyEmailUseCase {
	return &VerifyEmailUseCase{
		userRepo:              userRepo,
		emailVerificationRepo: emailVerificationRepo,
	}
}

func (uc *VerifyEmailUseCase) Execute(verificationToken string, verificationCode string) (bool, error) {
	ctx := context.TODO()

	userId, verificationExpiresAt, err := security.ValidateVerificationToken(verificationToken)
	if err != nil {
		return false, err
	}

	if time.Now().After(verificationExpiresAt) {
		return false, errors.New("verification token has expired")
	}

	emailVerification, err := uc.emailVerificationRepo.GetByUserID(ctx, userId)
	if err != nil {
		return false, err
	}

	if emailVerification.Code != verificationCode {
		return false, errors.New("verification code is incorrect")
	}

	err = uc.emailVerificationRepo.MarkAsVerified(ctx, userId)
	if err != nil {
		return false, err
	}

	return true, nil
}
