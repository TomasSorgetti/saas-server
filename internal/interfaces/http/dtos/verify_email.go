package dtos

type VerifyEmailInput struct {
	VerificationToken string `json:"verification_token" binding:"required"`
	VerificationCode  string `json:"verification_code" binding:"required"`
}