package auth

import (
	"luthierSaas/internal/infrastructure/cache"
	"luthierSaas/internal/infrastructure/email"
	"luthierSaas/internal/interfaces/repository"

	"github.com/rs/zerolog"
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
    Logout *LogoutUseCase
}

func NewAuthUseCases(
    userRepo repository.UserRepository, 
    suscriptionRepo repository.SubscriptionRepository, 
    emailVerificationRepo repository.EmailVerificationRepository, 
    sessionRepo repository.SessionRepository, 
    emailService *email.EmailService, 
    cacheService *cache.Cache,
    logger      *zerolog.Logger,
    ) *AuthUseCases{
        
    return &AuthUseCases{
        Login:      NewLoginUseCase(userRepo, emailVerificationRepo, sessionRepo, emailService, logger),
        Register:   NewRegisterUserUseCase(userRepo, suscriptionRepo, emailService, cacheService, logger),
        CheckEmail: NewCheckEmailUseCase(userRepo, cacheService),
        VerifyEmail: NewVerifyEmailUseCase(userRepo, emailVerificationRepo),
        ResendVerificationCode: NewResendVerificationCodeUseCase(userRepo, emailVerificationRepo, emailService),
        RefreshToken: NewRefreshTokenUseCase(userRepo, sessionRepo),
        GoogleLogin: NewLoginGoogleUseCase(userRepo, emailVerificationRepo, emailService),
        GoogleCallback: NewGoogleCallbackUseCase(),
        Logout: NewLogoutUseCase(sessionRepo, logger),
    }
}