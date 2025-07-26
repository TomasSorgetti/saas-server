package auth

import (
	"luthierSaas/internal/infrastructure/email"
	"luthierSaas/internal/interfaces/repository"
)

type AuthUseCases struct {
    Register   *RegisterUserUseCase
    CheckEmail *CheckEmailUseCase
    VerifyEmail *VerifyEmailUseCase
}

func NewAuthUseCases(userRepo repository.UserRepository, emailVerificationRepo repository.EmailVerificationRepository, emailService *email.EmailService) *AuthUseCases{
    return &AuthUseCases{
        Register:   NewRegisterUserUseCase(userRepo, emailService),
        CheckEmail: NewCheckEmailUseCase(userRepo),
        VerifyEmail: NewVerifyEmailUseCase(userRepo, emailVerificationRepo),
    }
}