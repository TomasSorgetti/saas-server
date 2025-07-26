package entities

import "time"

type EmailVerification struct {
	ID        int64
	UserID    int64
	Code      string
	ExpiresAt time.Time
	Verified  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
