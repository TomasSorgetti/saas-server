package user

import (
	"errors"
	"luthierSaas/internal/infrastructure/cache"
	"luthierSaas/internal/infrastructure/security"
	"luthierSaas/internal/interfaces/http/dtos"
	"luthierSaas/internal/interfaces/repository"
)

type ChangePasswordUseCase struct {
	userRepo repository.UserRepository
	cache    *cache.Cache
}

func NewChangePasswordUseCase(userRepo repository.UserRepository, cache *cache.Cache) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{userRepo, cache}
}

func (uc *ChangePasswordUseCase) Execute(userID int, input dtos.ChangePasswordInput) (error) {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil{
		return errors.New("user not found")
	}
    
    // no es necesario supongo
    // if user.LoginMethod != nil{
	// 	return errors.New("invalid login method")
	// }

    if !security.ComparePasswords(user.Password, input.Password) {
        return errors.New("invalid current password")
    }

    if security.ComparePasswords(user.Password, input.NewPassword) {
        return errors.New("new password cannot be the same as the current password")
    }

	hashedPassword, err := security.HashPassword(input.NewPassword)
    if err != nil {
        return errors.New("failed to hash new password")
    }

    err = uc.userRepo.UpdatePassword(userID, hashedPassword)
    if err != nil {
        return errors.New("failed to update password")
    }

    return nil
}