package domain

import "time"

type DailyStats struct {
	Date             time.Time
	TotalDuration    time.Duration
	MouseDuration    time.Duration
	KeyboardDuration time.Duration
} 