package auth

import (
	"context"
	"errors"
	"fmt"
	"luthierSaas/internal/domain/entities"
	"luthierSaas/internal/infrastructure/email"
	"luthierSaas/internal/infrastructure/security"
	"luthierSaas/internal/interfaces/http/dtos"
	"luthierSaas/internal/interfaces/repository"
	"time"

	"github.com/rs/zerolog"
)

type LoginUseCase struct {
	userRepo              repository.UserRepository
	emailVerificationRepo  repository.EmailVerificationRepository
	sessionRepo 			repository.SessionRepository
	emailService *email.EmailService
	logger      *zerolog.Logger
}


func NewLoginUseCase(
	userRepo repository.UserRepository, 
	emailVerificationRepo repository.EmailVerificationRepository, 
	sessionRepo repository.SessionRepository, 
	emailService *email.EmailService,
	logger *zerolog.Logger) *LoginUseCase {

	return &LoginUseCase{
		userRepo:     userRepo,
		emailVerificationRepo: emailVerificationRepo,
		sessionRepo: sessionRepo,
		emailService: emailService,
		logger: logger,
	}
}

func (uc *LoginUseCase) Execute(input dtos.LoginInput, deviceInfo string) (*dtos.LoginResponse, error) {
	user, err := uc.userRepo.FindByEmail(input.Email)
	if err != nil {
		uc.logger.Error().
            Err(err).
            Str("email", input.Email).
            Msg("Failed to find user by email")
		return nil, err
	}

	if user == nil {
		uc.logger.Error().
            Str("email", input.Email).
            Msg("User not found")
		return nil, errors.New("user not found")
	}

	if user.LoginMethod != nil {
		uc.logger.Error().
            Int("user_id", user.ID).
            Str("email", input.Email).
            Str("login_method", *user.LoginMethod).
            Msg("Invalid login method")
		return nil, errors.New("invalid login method")
	}

	if user.Deleted {
		uc.logger.Error().
            Int("user_id", user.ID).
            Str("email", input.Email).
            Msg("User deleted tried to login")
		return nil, errors.New("account deleted")
	}	

	if !security.ComparePasswords(user.Password, input.Password) {
		uc.logger.Error().
            Int("user_id", user.ID).
            Str("email", input.Email).
            Msg("Invalid credentials")
		return nil, errors.New("invalid credentials")
	}

	if !user.Verified {
        ctx := context.TODO()
        emailVerification, err := uc.emailVerificationRepo.GetByUserID(ctx, user.ID)
        if err != nil {
			uc.logger.Error().
                Err(err).
                Int("user_id", user.ID).
                Str("email", input.Email).
                Msg("Failed to get email verification by user id")
            return nil, err
        }

        var verificationToken, verificationCode string
		expiresAt := time.Now().Add(15 * time.Minute) 

		verificationToken, err = security.CreateVerificationToken(user.ID, user.Email, expiresAt)

		if err != nil {
			uc.logger.Error().
                Err(err).
                Int("user_id", user.ID).
                Str("email", input.Email).
                Msg("Failed to create verification token")
			return nil, err
		}
        if emailVerification == nil || time.Now().After(emailVerification.ExpiresAt) {
            verificationCode, err = security.GenerateVerificationCode(6)
			if err != nil {
				uc.logger.Error().
                    Err(err).
                    Int("user_id", user.ID).
                    Str("email", input.Email).
                    Msg("Failed to generate verification code")
                return nil, err
            }


            err = uc.emailVerificationRepo.UpdateCode(ctx, emailVerification.ID, verificationCode, expiresAt)
			
            if err != nil {
				uc.logger.Error().
                    Err(err).
                    Int("user_id", user.ID).
                    Str("email", input.Email).
                    Msg("Failed to update verification code")
                return nil, err
            }
        } else {
            verificationCode = emailVerification.Code
        }

        emailJob := email.EmailJob{
			To:      user.Email,
			Subject: "Verificá tu cuenta",
			Body:    fmt.Sprintf("Tu código de verificación es: %s", verificationCode),
		}

		if err := uc.emailService.SendEmailAsync(context.Background(), emailJob); err != nil {
			uc.logger.Error().
                Err(err).
                Int("user_id", user.ID).
                Str("email", input.Email).
                Msg("Failed to send verification email")
			return nil, fmt.Errorf("falló el envío del email de verificación: %w", err)
		}

		uc.logger.Info().
            Int("user_id", user.ID).
            Str("email", input.Email).
            Msg("Verification email sent successfully")

		return &dtos.LoginResponse{
			VerificationRequired: true,
			VerificationExpiresAt: expiresAt,
            VerificationToken:   verificationToken,
            Redirect:            "/verify",
        }, nil
    }

	// Update LastLogin timestamp
    currentTime := time.Now()
    err = uc.userRepo.UpdateLastLogin(context.TODO(), user.ID, currentTime)
    if err != nil {
		uc.logger.Error().
            Err(err).
            Int("user_id", user.ID).
            Str("email", input.Email).
            Msg("Failed to update last login")
        return nil, err
    }

	user.LastLogin = currentTime.Format(time.RFC3339)

	accessToken, err := security.CreateAccessToken(user.ID)
	if err != nil {
		uc.logger.Error().
            Err(err).
            Int("user_id", user.ID).
            Str("email", input.Email).
            Msg("Failed to create access token")
		return nil, err
	}

	refreshToken, err := security.CreateRefreshToken(user.ID)
	if err != nil {
		uc.logger.Error().
            Err(err).
            Int("user_id", user.ID).
            Str("email", input.Email).
            Msg("Failed to create refresh token")
		return nil, err
	}

	accessTokenHash, err := security.HashToken(accessToken)
    if err != nil {
		uc.logger.Error().
            Err(err).
            Int("user_id", user.ID).
            Str("email", input.Email).
            Msg("Failed to hash access token")
        return nil, fmt.Errorf("failed to hash access token: %w", err)
    }

    refreshTokenHash, err := security.HashToken(refreshToken)
    if err != nil {
		uc.logger.Error().
            Err(err).
            Int("user_id", user.ID).
            Str("email", input.Email).
            Msg("Failed to hash refresh token")
        return nil, fmt.Errorf("failed to hash refresh token: %w", err)
    }

    session := &entities.Session{
        UserID:           user.ID,
        AccessTokenHash:  accessTokenHash,
        RefreshTokenHash: refreshTokenHash,
        ExpiresAt:        time.Now().Add(15 * time.Minute), 
        RefreshExpiresAt: time.Now().Add(7 * 24 * time.Hour), 
        IsValid:          true,
		DeviceInfo:       deviceInfo,
    }

    ctx := context.TODO()
    err = uc.sessionRepo.Create(ctx, session)
    if err != nil {
		uc.logger.Error().
            Err(err).
            Int("user_id", user.ID).
            Str("email", input.Email).
            Str("device_info", deviceInfo).
            Msg("Failed to create session")
        return nil, fmt.Errorf("failed to create session: %w", err)
    }

	uc.logger.Info().
        Int("user_id", user.ID).
        Str("email", input.Email).
        Str("device_info", deviceInfo).
        Str("access_token_hash", accessTokenHash).
        Msg("User logged in successfully")
		
	profile := &dtos.ProfileResponse{
        ID:           user.ID,
        Email:        user.Email,
        FirstName:    user.FirstName,
        LastName:     user.LastName,
        Phone:        user.Phone,
        Address:      user.Address,
        Country:      user.Country,
        WorkshopName: user.WorkshopName,
        LastLogin:    user.LastLogin,
		LoginMethod:  func() string { if user.LoginMethod != nil { return *user.LoginMethod }; return "" }(),
		Subscription: user.Subscription,
    }

	return &dtos.LoginResponse{
		Profile: profile,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}