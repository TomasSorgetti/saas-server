package entities

import "time"

type Session struct {
    ID                int64
    UserID            int64
    AccessTokenHash   string
    RefreshTokenHash  string
    ExpiresAt         time.Time
    RefreshExpiresAt  time.Time
    IsValid           bool
    CreatedAt         time.Time
    UpdatedAt         time.Time
}