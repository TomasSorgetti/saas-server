package auth

import (
	"context"
	"fmt"
	"luthierSaas/internal/infrastructure/email"
	"luthierSaas/internal/infrastructure/security"
	"luthierSaas/internal/interfaces/repository"
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

func (uc *ResendVerificationCodeUseCase) Execute(verificationToken string) (bool, error) {
	userId, userEmail, _, err := security.ValidateVerificationToken(verificationToken)
	if err != nil {
		return false, err
	}

	emailVerificationCode, err := uc.emailVerificationRepo.GetByUserID(context.Background(), userId)
	if err != nil {
		return false, err
	}

	if emailVerificationCode == nil {
		return false, nil
	}

	newCode, err := security.GenerateVerificationCode(6)
	if err != nil {
		return false, err
	}

	err = uc.emailVerificationRepo.UpdateCode(context.Background(), emailVerificationCode.ID, newCode)
	if err != nil {
		return false, err
	}

	emailJob := email.EmailJob{
		To:      userEmail,
		Subject: "Verificá tu cuenta",
		Body:    fmt.Sprintf("Tu código de verificación es: %s", newCode),
	}

	if err := uc.emailService.SendEmailAsync(context.Background(), emailJob); err != nil {
		// Deberia loguear el error - NOT_IMPLEMENTED
		return false, fmt.Errorf("falló el envío del email de verificación: %w", err)
	}

	return true, nil
}
