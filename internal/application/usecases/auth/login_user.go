package auth

import (
	"context"
	"errors"
	"fmt"
	"luthierSaas/internal/infrastructure/email"
	"luthierSaas/internal/infrastructure/security"
	"luthierSaas/internal/interfaces/http/dtos"
	"luthierSaas/internal/interfaces/repository"
	"time"
)

type LoginUseCase struct {
	userRepo              repository.UserRepository
	emailVerificationRepo  repository.EmailVerificationRepository
	emailService *email.EmailService
}


func NewLoginUseCase(userRepo repository.UserRepository, emailVerificationRepo repository.EmailVerificationRepository, emailService *email.EmailService) *LoginUseCase {
	return &LoginUseCase{
		userRepo:     userRepo,
		emailVerificationRepo: emailVerificationRepo,
		emailService: emailService,
	}
}

func (uc *LoginUseCase) Execute(input dtos.LoginInput) (*dtos.LoginResponse, error) {
	user, err := uc.userRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	// NOT IMPLEMENTED: Login method handling
	loginMethod := "password"
	if (loginMethod != "password") {
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


            err = uc.emailVerificationRepo.UpdateCode(ctx, emailVerification.ID, verificationCode)
			
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

	accessToken, err := security.CreateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := security.CreateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	profile := &dtos.ProfileResponse{
        ID:           user.ID,
        Email:        user.Email,
        FirstName:    user.FirstName,
        LastName:     user.LastName,
        Phone:        user.Phone,
        Country:      user.Country,
        WorkshopName: user.WorkshopName,
        LastLogin:    user.LastLogin,
		Subscription: user.Subscription,
    }

	return &dtos.LoginResponse{
		Profile: profile,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}