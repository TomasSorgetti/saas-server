package auth

import (
	"context"
	"fmt"
	"luthierSaas/internal/domain/entities"
	"luthierSaas/internal/infrastructure/email"
	"luthierSaas/internal/infrastructure/security"
	"luthierSaas/internal/interfaces/http/dtos"
	"luthierSaas/internal/interfaces/repository"
	"time"
)

type RegisterUserUseCase struct {
	userRepo     repository.UserRepository
	emailService *email.EmailService
}


func NewRegisterUserUseCase(userRepo repository.UserRepository, emailService *email.EmailService) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userRepo:     userRepo,
		emailService: emailService,
	}
}

func (uc *RegisterUserUseCase) Execute(input dtos.RegisterInput) (*dtos.RegisterResponse, error) {
	hashedPassword, err := security.HashPassword(input.Password)
	if err != nil {
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
	}

	userID, err := uc.userRepo.Save(user)
	if err != nil {
		return nil, err
	}

	code, err := security.GenerateVerificationCode(6)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(24 * time.Hour)

	err = uc.userRepo.CreateEmailVerification(userID, code, expiresAt)
	if err != nil {
		return nil, err
	}

	verificationToken, err := security.CreateVerificationToken(userID, expiresAt)
	if err != nil {
		return nil, err
	}
	
	emailJob := email.EmailJob{
		To:      user.Email,
		Subject: "Verificá tu cuenta",
		Body:    fmt.Sprintf("Tu código de verificación es: %s", code),
	}

	if err := uc.emailService.SendEmailAsync(context.Background(), emailJob); err != nil {
		// Deberia loguear el error - NOT_IMPLEMENTED
		return nil, fmt.Errorf("falló el envío del email de verificación: %w", err)
	}

	return &dtos.RegisterResponse{
		VerificationToken: verificationToken,
		VerificationExpiresAt: expiresAt,
	}, nil
}