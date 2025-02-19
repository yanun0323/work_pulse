package main

import (
	"database/sql"
	"fmt"
	"log"
	"main/internal/domain"
	"main/internal/repository/sqlite"
	"main/internal/ui/window"
	"main/internal/usecase"

	"fyne.io/fyne/v2/app"
	_ "github.com/mattn/go-sqlite3"

	hook "github.com/robotn/gohook"
)

func main() {
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	settings := usecase.NewSettingsManager()
	if err := settings.LoadSettings(); err != nil {
		log.Printf("載入設定時發生錯誤: %v", err)
	}

	repo := sqlite.NewSQLiteActivityRepository(db)
	tracker := usecase.NewActivityTracker(repo)
	tracker.UpdateThreshold(settings.GetSettings().ThresholdSeconds)

	// 清理未完成的活動
	if err := tracker.CleanupUnfinishedActivities(); err != nil {
		log.Printf("清理未完成活動時發生錯誤: %v", err)
	}

	// 啟動閾值檢查器
	tracker.StartThresholdChecker()

	// 啟動監聽程序
	go trackActivity(tracker)

	myApp := app.New()
	mainWindow := window.NewMainWindow(myApp, tracker, settings)
	mainWindow.Show()
}

func initDB() (*sql.DB, error) {
	dbPath := "workpulse.db"
	log.Printf("正在連接資料庫: %s", dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("無法開啟資料庫: %v", err)
	}

	// 測試連接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("無法連接資料庫: %v", err)
	}

	log.Printf("資料庫連接成功")

	// 創建活動記錄表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS activities (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			start_time INTEGER NOT NULL,  -- Unix timestamp in seconds
			end_time INTEGER,            -- Unix timestamp in seconds, NULL means not ended
			activity_type TEXT NOT NULL
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("無法創建資料表: %v", err)
	}

	log.Printf("資料表檢查/創建完成")
	return db, nil
}

func trackActivity(tracker *usecase.ActivityTracker) {
	evChan := hook.Start()
	defer hook.End()

	for ev := range evChan {
		switch ev.Kind {
		case hook.MouseMove, hook.MouseDrag, hook.MouseDown, hook.MouseUp,
			hook.KeyDown, hook.KeyUp:

			activityType := domain.MouseActivity
			if ev.Kind == hook.KeyDown || ev.Kind == hook.KeyUp {
				activityType = domain.KeyboardActivity
			}

			if err := tracker.StartActivity(activityType); err != nil {
				log.Printf("開始活動時發生錯誤: %v", err)
			}
		}
	}
}
