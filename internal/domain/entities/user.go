package entities

import "time"

type User struct {
	ID           int
	Email        string
	Password     string
	GoogleID     string
	LoginMethod  *string
	Role         string
	FirstName    string
	LastName     string
	Phone        string
	Address      string
	Country      string
	WorkshopName string
	IsActive     bool
	Deleted      bool
	Verified     bool
	CreatedAt    time.Time
	LastLogin    string
	Subscription *Subscription
}