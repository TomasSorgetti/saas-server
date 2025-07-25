package auth

import (
	"luthierSaas/internal/domain/entities"
	"luthierSaas/internal/infrastructure/email"
	"luthierSaas/internal/interfaces/http/dtos"
	"luthierSaas/internal/interfaces/repository"
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

func (uc *RegisterUserUseCase) Execute(input dtos.RegisterInput) error {
	user := &entities.User{
		Email:                input.Email,
		Password:             hashPassword(input.Password),
		Role:                 "luthier",
		FirstName:            input.FirstName,
		LastName:             input.LastName,
		Phone:                input.Phone,
		Address:              input.Address,
		Country:              input.Country,
		WorkshopName:         input.WorkshopName,
		IsActive:             true,
		Deleted:              false,
		SubscriptionPlan:     "free",
		SubscriptionStatus:   "active",
		LastLogin:            "",
		ResetPasswordToken:   "",
		ResetPasswordExpires: "",
	}
	return uc.userRepo.Save(user)
}

func hashPassword(password string) string {
	return password 
}