package entities

type SubscriptionPlan struct {
	ID           int64
	Name         string
	Description  string
	Price        float64
	DurationDays int
	CreatedAt    string
	UpdatedAt    string
}