package component

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"

	"main/internal/domain"
	"main/internal/usecase"
)

type TimelineChart struct {
	widget.BaseWidget
	tracker   *usecase.ActivityTracker
	startTime time.Time
	endTime   time.Time
}

func NewTimelineChart(tracker *usecase.ActivityTracker) *TimelineChart {
	chart := &TimelineChart{
		tracker:   tracker,
		startTime: time.Now().Truncate(24 * time.Hour),
		endTime:   time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour),
	}
	chart.ExtendBaseWidget(chart)
	return chart
}

func (c *TimelineChart) CreateRenderer() fyne.WidgetRenderer {
	return &timelineRenderer{
		chart: c,
		bg:    canvas.NewRectangle(color.White),
	}
}

type timelineRenderer struct {
	chart *TimelineChart
	bg    *canvas.Rectangle
	rects []fyne.CanvasObject
}

func (r *timelineRenderer) MinSize() fyne.Size {
	return fyne.NewSize(600, 150)
}

func (r *timelineRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	r.drawTimeline(size)
}

func (r *timelineRenderer) Refresh() {
	r.drawTimeline(r.bg.Size())
	canvas.Refresh(r.chart)
}

func (r *timelineRenderer) Objects() []fyne.CanvasObject {
	objects := make([]fyne.CanvasObject, 0, len(r.rects)+1)
	objects = append(objects, r.bg)
	objects = append(objects, r.rects...)
	return objects
}

func (r *timelineRenderer) Destroy() {}

func (r *timelineRenderer) drawTimeline(size fyne.Size) {
	// 使用本地時間
	today := time.Now().Local().Truncate(24 * time.Hour)
	r.chart.startTime = today
	r.chart.endTime = today.Add(24 * time.Hour)

	// 清除舊的矩形
	r.rects = nil

	// 設置內部繪圖區域（考慮padding）
	padding := float32(20)
	innerSize := fyne.NewSize(size.Width-padding*2, size.Height-padding*2)
	innerPos := fyne.NewPos(padding, padding)

	// 繪製背景
	background := canvas.NewRectangle(color.White)
	background.Resize(innerSize)
	background.Move(innerPos)
	r.rects = append(r.rects, background)

	// 繪製背景時間刻度
	for i := 0; i < 24; i++ {
		// 垂直線
		line := canvas.NewLine(color.NRGBA{R: 200, G: 200, B: 200, A: 255})
		x := innerPos.X + float32(i)*innerSize.Width/24
		line.Position1 = fyne.NewPos(x, innerPos.Y)
		line.Position2 = fyne.NewPos(x, innerPos.Y+innerSize.Height-20)
		line.StrokeWidth = 1
		r.rects = append(r.rects, line)

		// 時間標籤
		label := canvas.NewText(fmt.Sprintf("%02d:00", i), color.NRGBA{R: 100, G: 100, B: 100, A: 255})
		label.TextSize = 10
		label.Move(fyne.NewPos(x-10, innerPos.Y+innerSize.Height-15))
		r.rects = append(r.rects, label)
	}

	activities, err := r.chart.tracker.GetTodayActivities()
	if err != nil {
		log.Printf("獲取今日活動失敗: %v", err)
		return
	}

	// 記錄並打印活動時間，用於調試
	log.Printf("今日活動數量: %d", len(activities))
	for _, activity := range activities {
		endTimeStr := "進行中"
		if activity.IsEnded() {
			endTimeStr = activity.EndTime().Format("15:04:05")
		}
		log.Printf("活動: 開始=%v, 結束=%v, 類型=%v",
			activity.StartTime().Format("15:04:05"),
			endTimeStr,
			activity.Type)
	}

	// 繪製活動時間條
	for _, activity := range activities {
		startX := r.timeToX(activity.StartTime(), innerSize.Width)
		var endX float32

		// 處理正在進行中的活動
		if !activity.IsEnded() {
			// 如果活動還在進行中，使用當前時間作為結束時間
			if r.chart.tracker.IsActive() {
				endX = r.timeToX(time.Now(), innerSize.Width)
			} else {
				// 如果活動已經停止但沒有結束時間，使用最後活動時間
				endX = r.timeToX(r.chart.tracker.GetLastActivityTime(), innerSize.Width)
			}
		} else {
			endX = r.timeToX(activity.EndTime(), innerSize.Width)
		}

		// 打印計算出的座標，用於調試
		log.Printf("時間條座標: startX=%.2f, endX=%.2f", startX, endX)

		if endX > startX { // 確保時間條有寬度
			rect := canvas.NewRectangle(getActivityColor(activity.Type))
			rect.Move(fyne.NewPos(
				innerPos.X+startX,
				innerPos.Y+5,
			))
			rect.Resize(fyne.NewSize(
				endX-startX,
				innerSize.Height-30,
			))

			// 添加半透明效果，進行中的活動使用不同的透明度
			rectColor := rect.FillColor.(color.NRGBA)
			if !activity.IsEnded() {
				rectColor.A = 120 // 進行中的活動更透明
			} else {
				rectColor.A = 180
			}
			rect.FillColor = rectColor

			r.rects = append(r.rects, rect)
		}
	}
}

func (r *timelineRenderer) timeToX(t time.Time, width float32) float32 {
	// 將時間轉換為本地時間
	t = t.In(time.Local)

	// 只取時分秒部分
	hour, min, sec := t.Clock()

	// 計算從當天開始到現在的秒數
	seconds := float64(hour*3600 + min*60 + sec)
	totalSeconds := float64(24 * 60 * 60)
	proportion := seconds / totalSeconds

	return width * float32(proportion)
}

func getActivityColor(activityType domain.ActivityType) color.Color {
	switch activityType {
	case domain.MouseActivity:
		return color.NRGBA{R: 46, G: 204, B: 113, A: 255} // 綠色
	case domain.KeyboardActivity:
		return color.NRGBA{R: 52, G: 152, B: 219, A: 255} // 藍色
	default:
		return color.Gray{Y: 128}
	}
}
