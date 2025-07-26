package repository

import (
	"luthierSaas/internal/domain/entities"
	"time"
)

type UserRepository interface {
    Save(user *entities.User) (int64, error)
    CreateEmailVerification(userID int64, code string, expiresAt time.Time) error
    FindByID(id int64) (*entities.User, error)
    FindByEmail(email string) (*entities.User, error)
    FindAll() ([]*entities.User, error)
    UpdateEmailVerified(userID int64, verified bool) error
}