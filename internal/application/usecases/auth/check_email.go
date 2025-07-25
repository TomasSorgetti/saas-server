package auth

import "luthierSaas/internal/interfaces/repository"

type CheckEmailUseCase struct {
	userRepo repository.UserRepository
}

func NewCheckEmailUseCase(userRepo repository.UserRepository) *CheckEmailUseCase {
	return &CheckEmailUseCase{userRepo}
}

func (uc *CheckEmailUseCase) Execute(email string) (bool, error) {
	user, err := uc.userRepo.FindByEmail(email)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}