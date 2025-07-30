package dtos

import "time"

type LoginResponse struct {
	Profile               *ProfileResponse `json:"profile,omitempty"`
	AccessToken           string           `json:"access_token,omitempty"`
	RefreshToken          string           `json:"refresh_token,omitempty"`
	VerificationRequired  bool             `json:"verificationRequired,omitempty"`
	VerificationToken     string           `json:"verificationToken,omitempty"`
	VerificationExpiresAt time.Time        `json:"verificationCodeExpiresAt,omitempty"`
	Redirect              string           `json:"redirect,omitempty"`
}
