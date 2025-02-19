package window

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"main/internal/ui/component"
	"main/internal/usecase"
	"main/pkg/utils"
	"time"
)

type MainWindow struct {
	window   fyne.Window
	app      fyne.App
	tracker  *usecase.ActivityTracker
	settings *usecase.SettingsManager
}

func NewMainWindow(app fyne.App, tracker *usecase.ActivityTracker, settings *usecase.SettingsManager) *MainWindow {
	window := app.NewWindow("Work Pulse")
	return &MainWindow{
		window:   window,
		app:      app,
		tracker:  tracker,
		settings: settings,
	}
}

func (w *MainWindow) Show() {
	// 創建表格
	table := widget.NewTable(
		func() (int, int) {
			return len(w.tracker.GetDailyStats()) + 1, 4 // +1 for header row
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("000000000000")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			if i.Row == 0 {
				headers := []string{"日期", "總時間", "滑鼠時間", "鍵盤時間"}
				label.SetText(headers[i.Col])
				// 設置表頭樣式
				label.TextStyle = fyne.TextStyle{Bold: true}
			} else {
				stats := w.tracker.GetDailyStats()
				if i.Row-1 < len(stats) {
					stat := stats[i.Row-1]
					switch i.Col {
					case 0:
						label.SetText(stat.Date.Format("2006-01-02"))
					case 1:
						label.SetText(utils.FormatDuration(stat.TotalDuration))
					case 2:
						label.SetText(utils.FormatDuration(stat.MouseDuration))
					case 3:
						label.SetText(utils.FormatDuration(stat.KeyboardDuration))
					}
				}
			}
		},
	)

	// 設置表格列寬
	table.SetColumnWidth(0, 120) // 日期列
	table.SetColumnWidth(1, 100) // 總時間列
	table.SetColumnWidth(2, 100) // 滑鼠時間列
	table.SetColumnWidth(3, 100) // 鍵盤時間列

	// 創建一個固定高度的滾動容器
	tableContainer := container.NewVScroll(table)
	tableContainer.SetMinSize(fyne.NewSize(500, 300)) // 設置最小大小

	// 創建時間軸圖表
	timeline := component.NewTimelineChart(w.tracker)

	// 添加設定按鈕
	settingsBtn := widget.NewButton("設定", func() {
		settingsWindow := NewSettingsWindow(w.app, w.settings, func() {
			// 當設定更新時，重新載入設定
			w.tracker.UpdateThreshold(w.settings.GetSettings().ThresholdSeconds)
		})
		settingsWindow.Show()
	})

	content := container.NewVBox(
		container.NewHBox(
			widget.NewLabel("工作時間追蹤"),
			settingsBtn,
		),
		tableContainer, // 使用包裝後的表格容器
		widget.NewLabel("今日活動時間軸"),
		timeline,
	)

	w.window.SetContent(content)
	w.window.Resize(fyne.NewSize(800, 600))

	// 定期更新 UI
	go func() {
		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			table.Refresh()
			timeline.Refresh()
		}
	}()

	w.window.ShowAndRun()
}

// ... 實現其他方法 ...
