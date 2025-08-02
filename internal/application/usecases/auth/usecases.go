package auth

import (
	"luthierSaas/internal/infrastructure/cache"
	"luthierSaas/internal/infrastructure/email"
	"luthierSaas/internal/interfaces/repository"
)

type AuthUseCases struct {
    Login      *LoginUseCase
    Register   *RegisterUserUseCase
    CheckEmail *CheckEmailUseCase
    VerifyEmail *VerifyEmailUseCase
    ResendVerificationCode *ResendVerificationCodeUseCase
    RefreshToken *RefreshTokenUseCase
    GoogleLogin *LoginGoogleUseCase
    GoogleCallback *GoogleCallbackUseCase
}

func NewAuthUseCases(userRepo repository.UserRepository, suscriptionRepo repository.SubscriptionRepository, emailVerificationRepo repository.EmailVerificationRepository, emailService *email.EmailService, cacheService *cache.Cache) *AuthUseCases{
    return &AuthUseCases{
        Login:      NewLoginUseCase(userRepo, emailVerificationRepo, emailService),
        Register:   NewRegisterUserUseCase(userRepo, suscriptionRepo, emailService, cacheService),
        CheckEmail: NewCheckEmailUseCase(userRepo, cacheService),
        VerifyEmail: NewVerifyEmailUseCase(userRepo, emailVerificationRepo),
        ResendVerificationCode: NewResendVerificationCodeUseCase(userRepo, emailVerificationRepo, emailService),
        RefreshToken: NewRefreshTokenUseCase(userRepo),
        GoogleLogin: NewLoginGoogleUseCase(userRepo, emailVerificationRepo, emailService),
        GoogleCallback: NewGoogleCallbackUseCase(),
    }
}