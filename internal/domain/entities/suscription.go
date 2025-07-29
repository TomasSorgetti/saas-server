package entities

import "time"

type Subscription struct {
	ID        int       
	UserID    int       
	PlanID    int       
	PlanName  string    
	Status    string    
	StartedAt time.Time 
	ExpiresAt time.Time 
}