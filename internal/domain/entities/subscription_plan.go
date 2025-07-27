package entities

type SubscriptionPlan struct {
	ID           int
	Name         string
	Description  string
	Price        float64
	DurationDays int
	CreatedAt    string
	UpdatedAt    string
}