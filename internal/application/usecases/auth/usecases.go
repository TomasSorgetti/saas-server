package auth

import (
	"luthierSaas/internal/infrastructure/email"
	"luthierSaas/internal/interfaces/repository"
)

type AuthUseCases struct {
    Login      *LoginUseCase
    Register   *RegisterUserUseCase
    CheckEmail *CheckEmailUseCase
    VerifyEmail *VerifyEmailUseCase
    ResendVerificationCode *ResendVerificationCodeUseCase
}

func NewAuthUseCases(userRepo repository.UserRepository, emailVerificationRepo repository.EmailVerificationRepository, emailService *email.EmailService) *AuthUseCases{
    return &AuthUseCases{
        Login:      NewLoginUseCase(userRepo, emailVerificationRepo),
        Register:   NewRegisterUserUseCase(userRepo, emailService),
        CheckEmail: NewCheckEmailUseCase(userRepo),
        VerifyEmail: NewVerifyEmailUseCase(userRepo, emailVerificationRepo),
        ResendVerificationCode: NewResendVerificationCodeUseCase(userRepo, emailVerificationRepo, emailService),
    }
}