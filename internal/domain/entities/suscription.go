package entities

type Subscription struct {
	ID        int
	UserID    int
	PlanID    int
	Status    string
	StartedAt string
	ExpiresAt string
	CreatedAt string
	UpdatedAt string
}