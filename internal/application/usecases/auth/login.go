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
)

type LoginUseCase struct {
	userRepo              repository.UserRepository
	emailVerificationRepo  repository.EmailVerificationRepository
	sessionRepo 			repository.SessionRepository
	emailService *email.EmailService
}


func NewLoginUseCase(userRepo repository.UserRepository, emailVerificationRepo repository.EmailVerificationRepository, sessionRepo repository.SessionRepository, emailService *email.EmailService) *LoginUseCase {
	return &LoginUseCase{
		userRepo:     userRepo,
		emailVerificationRepo: emailVerificationRepo,
		sessionRepo: sessionRepo,
		emailService: emailService,
	}
}

func (uc *LoginUseCase) Execute(input dtos.LoginInput, deviceInfo string) (*dtos.LoginResponse, error) {
	user, err := uc.userRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	if user.LoginMethod != nil {
		return nil, errors.New("invalid login method")
	}

	if user.Deleted {
		return nil, errors.New("account deleted")
	}	

	if !security.ComparePasswords(user.Password, input.Password) {
		return nil, errors.New("invalid credentials")
	}

	if !user.Verified {
        ctx := context.TODO()
        emailVerification, err := uc.emailVerificationRepo.GetByUserID(ctx, user.ID)
        if err != nil {
            return nil, err
        }

        var verificationToken, verificationCode string
		expiresAt := time.Now().Add(15 * time.Minute) 

		verificationToken, err = security.CreateVerificationToken(user.ID, user.Email, expiresAt)

		if err != nil {
			return nil, err
		}
        if emailVerification == nil || time.Now().After(emailVerification.ExpiresAt) {
            verificationCode, err = security.GenerateVerificationCode(6)
			if err != nil {
                return nil, err
            }


            err = uc.emailVerificationRepo.UpdateCode(ctx, emailVerification.ID, verificationCode, expiresAt)
			
            if err != nil {
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
			// Deberia loguear el error - NOT_IMPLEMENTED
			return nil, fmt.Errorf("falló el envío del email de verificación: %w", err)
		}

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
        return nil, err
    }

	user.LastLogin = currentTime.Format(time.RFC3339)

	accessToken, err := security.CreateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := security.CreateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	accessTokenHash, err := security.HashToken(accessToken)
    if err != nil {
        return nil, fmt.Errorf("failed to hash access token: %w", err)
    }

    refreshTokenHash, err := security.HashToken(refreshToken)
    if err != nil {
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
        return nil, fmt.Errorf("failed to create session: %w", err)
    }

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