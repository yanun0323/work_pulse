package sqlite

import (
	"context"
	"database/sql"
	"log"
	"time"

	"main/internal/domain"
	"main/internal/repository"
)

var _ repository.ActivityRepository = &SQLiteActivityRepository{}

type SQLiteActivityRepository struct {
	db *sql.DB
}

func NewSQLiteActivityRepository(db *sql.DB) *SQLiteActivityRepository {
	return &SQLiteActivityRepository{db: db}
}

func (r *SQLiteActivityRepository) Save(ctx context.Context, activity domain.Activity) error {
	var endTime interface{}
	if activity.EndTimeUnix > 0 {
		endTime = activity.EndTimeUnix
	}

	result, err := r.db.ExecContext(ctx, `
		INSERT INTO activities (start_time, end_time, activity_type)
		VALUES (?, ?, ?)
	`,
		activity.StartTimeUnix,
		endTime,
		activity.Type,
	)
	if err != nil {
		log.Printf("保存活動失敗: %v", err)
		return err
	}

	id, _ := result.LastInsertId()
	log.Printf("保存活動成功: ID=%d", id)
	return nil
}

func (r *SQLiteActivityRepository) UpdateEndTime(ctx context.Context, activity domain.Activity) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE activities 
		SET end_time = ? 
		WHERE start_time = ? AND activity_type = ?
	`,
		activity.EndTimeUnix,
		activity.StartTimeUnix,
		activity.Type,
	)
	return err
}

func (r *SQLiteActivityRepository) GetActivities(ctx context.Context) ([]domain.Activity, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, start_time, end_time, activity_type 
		FROM activities 
		ORDER BY start_time DESC
	`)
	if err != nil {
		log.Printf("查詢活動失敗: %v", err)
		return nil, err
	}
	defer rows.Close()

	var activities []domain.Activity
	for rows.Next() {
		var activity domain.Activity
		var startTime, endTime sql.NullInt64
		var activityType string

		if err := rows.Scan(&activity.ID, &startTime, &endTime, &activityType); err != nil {
			log.Printf("掃描活動資料失敗: %v", err)
			return nil, err
		}

		activity.StartTimeUnix = startTime.Int64
		if endTime.Valid {
			activity.EndTimeUnix = endTime.Int64
		}
		activity.Type = domain.ActivityType(activityType)

		activities = append(activities, activity)
	}

	log.Printf("查詢到 %d 條活動記錄", len(activities))
	return activities, rows.Err()
}

func (r *SQLiteActivityRepository) GetTodayActivities(ctx context.Context) ([]domain.Activity, error) {
	today := time.Now().Truncate(24 * time.Hour).Unix()
	tomorrow := today + 24*60*60

	log.Printf("查詢今日活動: 開始時間=%v, 結束時間=%v",
		time.Unix(today, 0).Format(time.RFC3339),
		time.Unix(tomorrow, 0).Format(time.RFC3339))

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, start_time, end_time, activity_type 
		FROM activities 
		WHERE start_time >= ? AND start_time < ?
		ORDER BY start_time ASC
	`, today, tomorrow)
	if err != nil {
		log.Printf("查詢失敗: %v", err)
		return nil, err
	}
	defer rows.Close()

	var activities []domain.Activity
	for rows.Next() {
		var activity domain.Activity
		var endTime sql.NullInt64

		if err := rows.Scan(&activity.ID, &activity.StartTimeUnix, &endTime, &activity.Type); err != nil {
			log.Printf("掃描資料失敗: %v", err)
			return nil, err
		}

		if endTime.Valid {
			activity.EndTimeUnix = endTime.Int64
		}

		activities = append(activities, activity)
	}

	log.Printf("成功查詢到 %d 筆活動記錄", len(activities))
	return activities, nil
}

func (r *SQLiteActivityRepository) CleanupUnfinishedActivities(ctx context.Context) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM activities 
		WHERE end_time IS NULL OR end_time = 0
	`)
	if err != nil {
		log.Printf("清理未完成活動失敗: %v", err)
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		log.Printf("獲取影響行數失敗: %v", err)
		return err
	}

	log.Printf("已清理 %d 筆未完成的活動記錄", count)
	return nil
}

// ... 實現其他方法 ...
