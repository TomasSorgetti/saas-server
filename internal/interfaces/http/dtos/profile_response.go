package dtos

import "luthierSaas/internal/domain/entities"

type ProfileResponse struct {
	ID           int               `json:"id"`
	Email        string            `json:"email"`
	FirstName    string            `json:"first_name"`
	LastName     string            `json:"last_name"`
	Phone        string            `json:"phone"`
	Address      string            `json:"address"`
	Country      string            `json:"country"`
	WorkshopName string            `json:"workshop_name"`
	LastLogin    string            `json:"last_login"`
	Subscription *entities.Subscription `json:"subscription,omitempty"`
}