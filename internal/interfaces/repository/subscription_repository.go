package repository

import (
	"luthierSaas/internal/domain/entities"
)

type SubscriptionRepository interface {
	Save(subscription *entities.Subscription) (int, error)
	GetFreeTierPlanID()(int, error)
	GetFreeTierPlan()(*entities.SubscriptionPlan, error)
}