package entities

type Subscription struct {
	ID        int64
	UserID    int64
	PlanID    int64
	Status    string
	StartedAt string
	ExpiresAt string
	CreatedAt string
	UpdatedAt string
}