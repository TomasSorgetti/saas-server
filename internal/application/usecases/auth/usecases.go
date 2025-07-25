package auth

import (
	"luthierSaas/internal/infrastructure/email"
	"luthierSaas/internal/interfaces/repository"
)

type AuthUseCases struct {
    Register   *RegisterUserUseCase
    CheckEmail *CheckEmailUseCase
}

func NewAuthUseCases(userRepo repository.UserRepository, emailService *email.EmailService) *AuthUseCases{
    return &AuthUseCases{
        Register:   NewRegisterUserUseCase(userRepo, emailService),
        CheckEmail: NewCheckEmailUseCase(userRepo),
    }
}