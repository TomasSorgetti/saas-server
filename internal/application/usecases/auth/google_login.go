package auth

import (
	"fmt"
	"luthierSaas/internal/infrastructure/email"
	"luthierSaas/internal/interfaces/repository"
)

type LoginGoogleUseCase struct {
	userRepo              repository.UserRepository
	emailVerificationRepo  repository.EmailVerificationRepository
	emailService *email.EmailService
}


func NewLoginGoogleUseCase(userRepo repository.UserRepository, emailVerificationRepo repository.EmailVerificationRepository, emailService *email.EmailService) *LoginGoogleUseCase {
	return &LoginGoogleUseCase{
		userRepo:     userRepo,
		emailVerificationRepo: emailVerificationRepo,
		emailService: emailService,
	}
}

func (uc *LoginGoogleUseCase) Execute()  error {
	return fmt.Errorf("not implemented")
}