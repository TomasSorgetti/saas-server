package auth

import (
	"errors"
	"luthierSaas/internal/infrastructure/security"
	"luthierSaas/internal/interfaces/http/dtos"
	"luthierSaas/internal/interfaces/repository"
)

type LoginUseCase struct {
	userRepo              repository.UserRepository
	emailVerificationRepo  repository.EmailVerificationRepository
}


func NewLoginUseCase(userRepo repository.UserRepository, emailVerificationRepo repository.EmailVerificationRepository) *LoginUseCase {
	return &LoginUseCase{
		userRepo:     userRepo,
		emailVerificationRepo: emailVerificationRepo,
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
		return nil, errors.New("email not verified")
	}

	accessToken, err := security.CreateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := security.CreateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Devolver respuesta
	return &dtos.LoginResponse{
		UserID:       user.ID,
		Email:        user.Email,
		Role:         user.Role,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}