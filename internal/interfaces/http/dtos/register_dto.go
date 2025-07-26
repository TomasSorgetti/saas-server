package dtos

type RegisterInput struct {
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required"`
	FirstName    string `json:"firstName" binding:"required"`
	LastName     string `json:"lastName" binding:"required"`
	Phone        string `json:"phone"`
	Address      string `json:"address"`
	Country      string `json:"country"`
	WorkshopName string `json:"workshopName"`
}