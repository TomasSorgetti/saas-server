package repository

import "luthierSaas/internal/domain/entities"

type UserRepository interface {
    Save(user *entities.User) error
    FindByID(id int64) (*entities.User, error)
    FindByEmail(email string) (*entities.User, error)
    FindAll() ([]*entities.User, error)
}