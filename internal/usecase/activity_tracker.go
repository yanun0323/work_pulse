package usecase

import (
	"context"
	"log"
	"sort"
	"time"

	"main/internal/domain"
	"main/internal/repository"
)

type ActivityTracker struct {
	isActive         bool
	lastActivity     time.Time
	repo             repository.ActivityRepository
	activities       []domain.Activity
	thresholdSeconds int
}

func NewActivityTracker(repo repository.ActivityRepository) *ActivityTracker {
	return &ActivityTracker{
		repo:             repo,
		isActive:         false,
		activities:       make([]domain.Activity, 0),
		thresholdSeconds: 15, // 預設值
	}
}

func (t *ActivityTracker) StartActivity(activityType domain.ActivityType) error {
	if !t.isActive {
		activity := domain.Activity{
			Type: activityType,
		}
		activity.SetStartTime(time.Now())

		log.Printf("開始新活動: 類型=%v, 開始時間=%v", activityType, activity.StartTime())

		if err := t.repo.Save(context.Background(), activity); err != nil {
			return err
		}

		t.activities = append(t.activities, activity)
		t.isActive = true
	}
	t.lastActivity = time.Now()
	return nil
}

func (t *ActivityTracker) StopActivity() error {
	if t.isActive && len(t.activities) > 0 {
		lastIdx := len(t.activities) - 1
		t.activities[lastIdx].SetEndTime(t.lastActivity)

		log.Printf("停止活動: 結束時間=%v", t.activities[lastIdx].EndTime())

		if err := t.repo.UpdateEndTime(context.Background(), t.activities[lastIdx]); err != nil {
			return err
		}

		t.isActive = false
	}
	return nil
}

func (t *ActivityTracker) GetDailyStats() []domain.DailyStats {
	statsMap := make(map[string]*domain.DailyStats)

	for _, activity := range t.activities {
		startTime := activity.StartTime()
		// 跳過無效的時間
		if startTime.Year() < 2000 {
			continue
		}

		date := startTime.Format("2006-01-02")
		if _, exists := statsMap[date]; !exists {
			statsMap[date] = &domain.DailyStats{
				Date: startTime.Truncate(24 * time.Hour), // 確保日期是當天的開始時間
			}
		}

		var duration time.Duration
		if !activity.IsEnded() {
			if t.isActive {
				duration = time.Since(startTime)
			} else {
				duration = t.lastActivity.Sub(startTime)
			}
		} else {
			duration = activity.EndTime().Sub(startTime)
		}

		// 確保持續時間是有效的
		if duration < 0 {
			continue
		}

		stats := statsMap[date]
		stats.TotalDuration += duration
		if activity.Type == domain.MouseActivity {
			stats.MouseDuration += duration
		} else {
			stats.KeyboardDuration += duration
		}
	}

	var result []domain.DailyStats
	for _, stats := range statsMap {
		result = append(result, *stats)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Date.After(result[j].Date) // 改為降序排序，最新的日期在前
	})

	return result
}

func (t *ActivityTracker) LoadActivities() error {
	activities, err := t.repo.GetActivities(context.Background())
	if err != nil {
		return err
	}
	t.activities = activities
	log.Printf("載入活動記錄: 數量=%d", len(activities))
	return nil
}

func (t *ActivityTracker) GetTodayActivities() ([]domain.Activity, error) {
	return t.repo.GetTodayActivities(context.Background())
}

func (t *ActivityTracker) StartThresholdChecker() {
	ticker := time.NewTicker(time.Second)
	go func() {
		for range ticker.C {
			if t.isActive && time.Since(t.lastActivity) > time.Second*time.Duration(t.thresholdSeconds) {
				if err := t.StopActivity(); err != nil {
					log.Printf("停止活動時發生錯誤: %v", err)
				}
			}
		}
	}()
}

func (t *ActivityTracker) UpdateThreshold(seconds int) {
	t.thresholdSeconds = seconds
}

func (t *ActivityTracker) IsActive() bool {
	return t.isActive
}

func (t *ActivityTracker) GetLastActivityTime() time.Time {
	return t.lastActivity
}

func (t *ActivityTracker) CleanupUnfinishedActivities() error {
	if err := t.repo.CleanupUnfinishedActivities(context.Background()); err != nil {
		return err
	}
	return t.LoadActivities() // 重新載入活動列表
}

// ... 實現其他方法 ...
