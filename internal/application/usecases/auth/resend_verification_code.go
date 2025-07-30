package auth

import (
	"context"
	"fmt"
	"luthierSaas/internal/infrastructure/email"
	"luthierSaas/internal/infrastructure/security"
	"luthierSaas/internal/interfaces/http/dtos"
	"luthierSaas/internal/interfaces/repository"
	"time"
)

type ResendVerificationCodeUseCase struct {
	userRepo              repository.UserRepository
	emailVerificationRepo repository.EmailVerificationRepository
	emailService          *email.EmailService
}
func NewResendVerificationCodeUseCase(userRepo repository.UserRepository, emailVerificationRepo repository.EmailVerificationRepository, emailService *email.EmailService) *ResendVerificationCodeUseCase {
	return &ResendVerificationCodeUseCase{
		userRepo:              userRepo,
		emailVerificationRepo: emailVerificationRepo,
		emailService:          emailService,
	}
}

func (uc *ResendVerificationCodeUseCase) Execute(verificationToken string) (*dtos.ResendCode, error) {
	userId, userEmail, _, err := security.ValidateVerificationToken(verificationToken)
	if err != nil {
		return nil, err
	}

	emailVerificationCode, err := uc.emailVerificationRepo.GetByUserID(context.Background(), userId)
	if err != nil {
		return nil, err
	}

	if emailVerificationCode == nil {
		return nil, err
	}
	
	newCode, err := security.GenerateVerificationCode(6)
	if err != nil {
		return nil, err
	}
	
	newExpiresAt := time.Now().Add(15 * time.Minute)

	newVerificationToken, err := security.CreateVerificationToken(userId, userEmail, newExpiresAt)
	if err != nil {
		return nil, err
	}
	
	err = uc.emailVerificationRepo.UpdateCode(context.Background(), emailVerificationCode.ID, newCode, newExpiresAt)
	if err != nil {
		return nil, err
	}

	emailJob := email.EmailJob{
		To:      userEmail,
		Subject: "Verificá tu cuenta",
		Body:    fmt.Sprintf("Tu código de verificación es: %s", newCode),
	}

	if err := uc.emailService.SendEmailAsync(context.Background(), emailJob); err != nil {
		// Deberia loguear el error - NOT_IMPLEMENTED
		return nil, fmt.Errorf("falló el envío del email de verificación: %w", err)
	}

	return &dtos.ResendCode{
		VerificationExpiresAt: newExpiresAt,
        VerificationToken:   newVerificationToken,
	},nil
}
