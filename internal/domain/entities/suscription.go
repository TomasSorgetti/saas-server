package entities

import "time"

type Subscription struct {
	ID        int    		`json:"id"`
	UserID    int    		`json:"user_id"`   
    PlanID    int    		`json:"plan_id"`
    PlanName  string 		`json:"plan_name"`
    Status    string 		`json:"status"`
    StartedAt time.Time  	`json:"started_at"`
    ExpiresAt time.Time  	`json:"expires_at"`
}

