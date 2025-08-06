package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"luthierSaas/internal/domain/entities"
	"luthierSaas/internal/infrastructure/cache"
	"luthierSaas/internal/infrastructure/email"
	"luthierSaas/internal/infrastructure/security"
	"luthierSaas/internal/interfaces/http/dtos"
	"luthierSaas/internal/interfaces/repository"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
)

type GoogleCallbackUseCase struct {
	oauthConfig           *oauth2.Config
	userRepo              repository.UserRepository
	subscriptionRepo      repository.SubscriptionRepository
	emailVerificationRepo repository.EmailVerificationRepository
	sessionRepo           repository.SessionRepository
	emailService          *email.EmailService
	cacheService          *cache.Cache
	logger                *zerolog.Logger
}

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
}

func NewGoogleCallbackUseCase(
	oauthConfig *oauth2.Config,
	userRepo repository.UserRepository,
	subscriptionRepo repository.SubscriptionRepository,
	emailVerificationRepo repository.EmailVerificationRepository,
	sessionRepo repository.SessionRepository,
	emailService *email.EmailService,
	cacheService *cache.Cache,
	logger *zerolog.Logger,
) *GoogleCallbackUseCase {
	return &GoogleCallbackUseCase{
		oauthConfig:           oauthConfig,
		userRepo:              userRepo,
		subscriptionRepo:      subscriptionRepo,
		emailVerificationRepo: emailVerificationRepo,
		sessionRepo:           sessionRepo,
		emailService:          emailService,
		cacheService:          cacheService,
		logger:                logger,
	}
}

func (uc *GoogleCallbackUseCase) Execute(ctx context.Context, code, state, deviceInfo string) (*dtos.LoginResponse, error) {
	// Validar el estado CSRF
	cachedState, err := uc.cacheService.Get(ctx, "oauth_state:"+state)
	if err != nil || cachedState != state {
		uc.logger.Error().
			Str("state", state).
			Msg("Invalid or expired CSRF state")
		return nil, fmt.Errorf("invalid or expired CSRF state")
	}
	// Eliminar el estado del cache
	_ = uc.cacheService.Delete(ctx, "oauth_state:"+state)

	// Intercambiar el código por un token
	token, err := uc.oauthConfig.Exchange(ctx, code)
	if err != nil {
		uc.logger.Error().
			Err(err).
			Msg("Failed to exchange code")
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Obtener información del usuario
	client := uc.oauthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		uc.logger.Error().
			Err(err).
			Msg("Failed to get user info")
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	// Usar io.ReadAll en lugar de ioutil.ReadAll
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		uc.logger.Error().
			Err(err).
			Msg("Failed to read response")
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var googleUser GoogleUser
	if err := json.Unmarshal(body, &googleUser); err != nil {
		uc.logger.Error().
			Err(err).
			Msg("Failed to parse user info")
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	// Buscar o crear usuario
	user, err := uc.userRepo.FindByEmail(googleUser.Email)
	if err != nil {
		uc.logger.Error().
			Err(err).
			Str("email", googleUser.Email).
			Msg("Failed to find user by email")
		return nil, err
	}

	var userID int
	if user == nil {
		// Crear nuevo usuario
		loginMethod := "google"
		user = &entities.User{
			Email:       googleUser.Email,
			GoogleID:    googleUser.ID,
			FirstName:   googleUser.Name,
			Role:        "user",
			IsActive:    true,
			Deleted:     false,
			Verified:    true, 
			LoginMethod: &loginMethod,
			CreatedAt:   time.Now(),
		}
		userID, err = uc.userRepo.Save(user)
		if err != nil {
			uc.logger.Error().
				Err(err).
				Str("email", googleUser.Email).
				Msg("Failed to save user")
			return nil, fmt.Errorf("failed to save user: %w", err)
		}

		// Crear suscripción gratuita
		plan, err := uc.subscriptionRepo.GetFreeTierPlan()
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
			PlanID:    plan.ID,
			PlanName:    plan.Name,
			Status:    "active",
			StartedAt: now,
			ExpiresAt: now.Add(14 * 24 * time.Hour),
		}

		_, err = uc.subscriptionRepo.Save(subscription)
		if err != nil {
			uc.logger.Error().
				Err(err).
				Int("user_id", userID).
				Msg("Failed to create subscription")
			return nil, fmt.Errorf("failed to create subscription: %w", err)
		}
		user.Subscription = subscription
	} else {
		userID = user.ID
		if user.Deleted {
			uc.logger.Error().
				Int("user_id", userID).
				Str("email", googleUser.Email).
				Msg("User deleted tried to login")
			return nil, fmt.Errorf("account deleted")
		}
		if user.LoginMethod == nil || *user.LoginMethod != "google" {
			uc.logger.Error().
				Int("user_id", userID).
				Str("email", googleUser.Email).
				Str("login_method", func() string { if user.LoginMethod != nil { return *user.LoginMethod } else { return "" } }()).
				Msg("Invalid login method")
			return nil, fmt.Errorf("invalid login method")
		}
	}

	// Manejar verificación de email
	if !user.Verified {
		verificationCode, err := security.GenerateVerificationCode(6)
		if err != nil {
			uc.logger.Error().
				Err(err).
				Int("user_id", userID).
				Msg("Failed to generate verification code")
			return nil, fmt.Errorf("failed to generate verification code: %w", err)
		}

		expiresAt := time.Now().Add(15 * time.Minute)
		err = uc.userRepo.CreateEmailVerification(userID, verificationCode, expiresAt)
		if err != nil {
			uc.logger.Error().
				Err(err).
				Int("user_id", userID).
				Msg("Failed to create email verification")
			return nil, fmt.Errorf("failed to create email verification: %w", err)
		}

		verificationToken, err := security.CreateVerificationToken(userID, googleUser.Email, expiresAt)
		if err != nil {
			uc.logger.Error().
				Err(err).
				Int("user_id", userID).
				Msg("Failed to create verification token")
			return nil, fmt.Errorf("failed to create verification token: %w", err)
		}

		emailJob := email.EmailJob{
			To:      googleUser.Email,
			Subject: "Verificá tu cuenta",
			Body:    fmt.Sprintf("Tu código de verificación es: %s", verificationCode),
		}

		if err := uc.emailService.SendEmailAsync(ctx, emailJob); err != nil {
			uc.logger.Error().
				Err(err).
				Int("user_id", userID).
				Str("email", googleUser.Email).
				Msg("Failed to send verification email")
			return nil, fmt.Errorf("failed to send verification email: %w", err)
		}

		uc.logger.Info().
			Int("user_id", userID).
			Str("email", googleUser.Email).
			Msg("Verification email sent successfully")

		return &dtos.LoginResponse{
			VerificationRequired:  true,
			VerificationExpiresAt: expiresAt,
			VerificationToken:    verificationToken,
			Redirect:             "/verify",
		}, nil
	}

	// Actualizar LastLogin
	currentTime := time.Now()
	err = uc.userRepo.UpdateLastLogin(ctx, userID, currentTime)
	if err != nil {
		uc.logger.Error().
			Err(err).
			Int("user_id", userID).
			Msg("Failed to update last login")
		return nil, fmt.Errorf("failed to update last login: %w", err)
	}

	// Crear sesión y tokens
	accessToken, err := security.CreateAccessToken(userID)
	if err != nil {
		uc.logger.Error().
			Err(err).
			Int("user_id", userID).
			Msg("Failed to create access token")
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	refreshToken, err := security.CreateRefreshToken(userID)
	if err != nil {
		uc.logger.Error().
			Err(err).
			Int("user_id", userID).
			Msg("Failed to create refresh token")
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	accessTokenHash, err := security.HashToken(accessToken)
	if err != nil {
		uc.logger.Error().
			Err(err).
			Int("user_id", userID).
			Msg("Failed to hash access token")
		return nil, fmt.Errorf("failed to hash access token: %w", err)
	}

	refreshTokenHash, err := security.HashToken(refreshToken)
	if err != nil {
		uc.logger.Error().
			Err(err).
			Int("user_id", userID).
			Msg("Failed to hash refresh token")
		return nil, fmt.Errorf("failed to hash refresh token: %w", err)
	}

	session := &entities.Session{
		UserID:           userID,
		AccessTokenHash:  accessTokenHash,
		RefreshTokenHash: refreshTokenHash,
		ExpiresAt:        time.Now().Add(15 * time.Minute),
		RefreshExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		IsValid:          true,
		DeviceInfo:       deviceInfo,
	}

	err = uc.sessionRepo.Create(ctx, session)
	if err != nil {
		uc.logger.Error().
			Err(err).
			Int("user_id", userID).
			Str("device_info", deviceInfo).
			Msg("Failed to create session")
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	uc.logger.Info().
		Int("user_id", userID).
		Str("email", googleUser.Email).
		Str("device_info", deviceInfo).
		Msg("User logged in successfully via Google")

	profile := &dtos.ProfileResponse{
		ID:           userID,
		Email:        googleUser.Email,
		FirstName:    user.FirstName,
		LoginMethod:  "google",
		Subscription: user.Subscription,
	}

	return &dtos.LoginResponse{
		Profile:      profile,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}