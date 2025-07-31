package repository

import (
	"luthierSaas/internal/domain/entities"
	"time"
)

type UserRepository interface {
    Save(user *entities.User) (int, error)
    CreateEmailVerification(userID int, code string, expiresAt time.Time) error
    FindByID(id int) (*entities.User, error)
    FindByEmail(email string) (*entities.User, error)
    FindAll() ([]*entities.User, error)
    UpdateEmailVerified(userID int, verified bool) error
    EmailExists(email string) (bool, error)
}