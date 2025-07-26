package dtos

type VerifyEmailResendInput struct {
	VerificationToken string `json:"verification_token" binding:"required"`
}