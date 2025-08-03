package auth

import (
	"context"
	"fmt"
	"luthierSaas/internal/domain/entities"
	"luthierSaas/internal/infrastructure/cache"
	"luthierSaas/internal/infrastructure/email"
	"luthierSaas/internal/infrastructure/security"
	"luthierSaas/internal/interfaces/http/dtos"
	"luthierSaas/internal/interfaces/repository"
	"time"

	"github.com/rs/zerolog"
)

type RegisterUserUseCase struct {
	userRepo     repository.UserRepository
	subscriptionRepo     repository.SubscriptionRepository
	emailService *email.EmailService
	cacheService *cache.Cache
	logger      *zerolog.Logger
}


func NewRegisterUserUseCase(userRepo repository.UserRepository, subscriptionRepo repository.SubscriptionRepository, emailService *email.EmailService, cacheService *cache.Cache, logger *zerolog.Logger) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userRepo:     userRepo,
		subscriptionRepo: subscriptionRepo,
		emailService: emailService,
		cacheService: cacheService,
		logger: logger,
	}
}

func (uc *RegisterUserUseCase) Execute(ctx context.Context, input dtos.RegisterInput) (*dtos.RegisterResponse, error) {
	cacheKey := fmt.Sprintf("email:check:%s", input.Email)
    _ = uc.cacheService.Delete(ctx, cacheKey)

	hashedPassword, err := security.HashPassword(input.Password)
	if err != nil {
		uc.logger.Error().
            Err(err).
            Str("email", input.Email).
            Msg("Failed to hash password")
		return nil, err
	}

	user := &entities.User{
		Email:        input.Email,
		Password:     hashedPassword,
		Role:         "user",
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		Phone:        input.Phone,
		Address:      input.Address,
		Country:      input.Country,
		WorkshopName: input.WorkshopName,
		IsActive:     true,
		Deleted:      false,
		LastLogin:    "",
		Verified:     false,
	}

	userID, err := uc.userRepo.Save(user)
	if err != nil {
		uc.logger.Error().
            Err(err).
            Str("email", input.Email).
            Msg("Failed to save user")
		return nil, err
	}

	planID, err := uc.subscriptionRepo.GetFreeTierPlanID()
    if err != nil {
		uc.logger.Error().
            Err(err).
            Int("user_id", userID).
            Msg("Failed to get free tier plan ID")
        return nil, fmt.Errorf("failed to get Free Tier plan ID: %w", err)
    }

	now := time.Now()
    subscription := &entities.Subscription{
        UserID:    userID,
        PlanID:    planID,
        Status:    "active",
        StartedAt: now,
        ExpiresAt: now.Add(14 * 24 * time.Hour), // 14 days
    }

    _, err = uc.subscriptionRepo.Save(subscription)
    if err != nil {
		uc.logger.Error().
            Err(err).
            Int("user_id", userID).
            Int("plan_id", planID).
            Msg("Failed to create subscription")
        return nil, fmt.Errorf("failed to create subscription for user %d: %w", userID, err)
    }
	
	code, err := security.GenerateVerificationCode(6)
	if err != nil {
		uc.logger.Error().
            Err(err).
            Int("user_id", userID).
            Msg("Failed to generate verification code")
		return nil, err
	}

	expiresAt := time.Now().Add(15 * time.Minute)

	err = uc.userRepo.CreateEmailVerification(userID, code, expiresAt)
	if err != nil {
		uc.logger.Error().
            Err(err).
            Int("user_id", userID).
            Msg("Failed to create email verification")
		return nil, err
	}

	verificationToken, err := security.CreateVerificationToken(userID, user.Email, expiresAt)
	if err != nil {
		uc.logger.Error().
            Err(err).
            Int("user_id", userID).
            Msg("Failed to create verification token")
		return nil, err
	}
	
	emailJob := email.EmailJob{
		To:      user.Email,
		Subject: "Verificá tu cuenta",
		Body:    fmt.Sprintf("Tu código de verificación es: %s", code),
	}

	if err := uc.emailService.SendEmailAsync(context.Background(), emailJob); err != nil {
		uc.logger.Error().
            Err(err).
            Int("user_id", userID).
            Str("email", user.Email).
            Msg("Failed to send verification email")
		return nil, fmt.Errorf("falló el envío del email de verificación: %w", err)
	}

	uc.logger.Info().
        Int("user_id", userID).
        Str("email", user.Email).
        Msg("User registered successfully")

	return &dtos.RegisterResponse{
		VerificationToken: verificationToken,
		VerificationExpiresAt: expiresAt,
	}, nil
}