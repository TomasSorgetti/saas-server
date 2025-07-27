package entities

import "time"

type EmailVerification struct {
	ID        int
	UserID    int
	Code      string
	ExpiresAt time.Time
	Verified  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
