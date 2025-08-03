package entities

import "time"

type Session struct {
    ID                int
    UserID            int
    AccessTokenHash   string
    RefreshTokenHash  string
    ExpiresAt         time.Time
    RefreshExpiresAt  time.Time
    IsValid           bool
    DeviceInfo       string
    CreatedAt         time.Time
    UpdatedAt         time.Time
}