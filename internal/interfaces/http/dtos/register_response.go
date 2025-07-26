package dtos

import "time"

type RegisterResponse struct {
	VerificationToken     string    `json:"verificationToken"`
	VerificationExpiresAt time.Time `json:"verificationCodeExpiresAt"`
}