package repository

import (
	"context"

	"main/internal/domain"
)

type ActivityRepository interface {
	Save(ctx context.Context, activity domain.Activity) error
	UpdateEndTime(ctx context.Context, activity domain.Activity) error
	GetActivities(ctx context.Context) ([]domain.Activity, error)
	GetTodayActivities(ctx context.Context) ([]domain.Activity, error)
	CleanupUnfinishedActivities(ctx context.Context) error
}
