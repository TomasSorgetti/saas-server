package dtos

type CheckEmailInput struct {
	Email string `json:"email" binding:"required,email"`
}