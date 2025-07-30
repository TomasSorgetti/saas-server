package dtos

import "time"

type ResendCode struct {
	VerificationToken     string    `json:"verificationToken"`
	VerificationExpiresAt time.Time `json:"verificationCodeExpiresAt"`
}