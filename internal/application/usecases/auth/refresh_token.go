package auth

import (
	"errors"
	"luthierSaas/internal/infrastructure/security"
	"luthierSaas/internal/interfaces/http/dtos"
	"luthierSaas/internal/interfaces/repository"
)

type RefreshTokenUseCase struct {
	userRepo repository.UserRepository
}

func NewRefreshTokenUseCase(userRepo repository.UserRepository) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{userRepo}
}

func (uc *RefreshTokenUseCase) Execute(userID int) (*dtos.RefreshResponse, error) {
	accessToken, err := security.CreateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := security.CreateRefreshToken(userID)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepo.FindByID(userID)
    if err != nil {
        return nil, errors.New("user not found")
    }

	if user.Deleted {
        return nil, errors.New("user deleted")
    }

   if !user.Verified {
    return nil, errors.New("user not verified")
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
    }

	
	return &dtos.RefreshResponse{
		Profile : profile,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}