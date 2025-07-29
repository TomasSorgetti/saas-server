package repositories

import (
	"database/sql"
	"fmt"
	"luthierSaas/internal/domain/entities"
)

type SubscriptionRepository struct {
	db *sql.DB
}


func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Save(subscription *entities.Subscription) (int, error) {
    query := `
        INSERT INTO subscriptions (user_id, plan_id, status, started_at, expires_at)
        VALUES (?, ?, ?, ?, ?)
    `
    result, err := r.db.Exec(query,
        subscription.UserID,
        subscription.PlanID,
        subscription.Status,
        subscription.StartedAt.Format("2006-01-02 15:04:05"),
        subscription.ExpiresAt.Format("2006-01-02 15:04:05"),
    )
    if err != nil {
        return 0, fmt.Errorf("failed to save subscription for user %d: %w", subscription.UserID, err)
    }

    id, err := result.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("failed to get subscription ID: %w", err)
    }

    subscription.ID = int(id)
    return int(id), nil
}

func (r *SubscriptionRepository) GetFreeTierPlanID() (int, error) {
    query := `SELECT id FROM subscription_plans WHERE name = 'Free Tier'`
    var planID int
    err := r.db.QueryRow(query).Scan(&planID)
    if err == sql.ErrNoRows {
        return 0, fmt.Errorf("free tier plan not found")
    }
    if err != nil {
        return 0, fmt.Errorf("failed to query Free Tier plan ID: %w", err)
    }
    return planID, nil
}