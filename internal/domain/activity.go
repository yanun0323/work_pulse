package domain

import "time"

type ActivityType string

const (
	MouseActivity    ActivityType = "mouse"
	KeyboardActivity ActivityType = "keyboard"
)

type Activity struct {
	ID            int64
	StartTimeUnix int64 // Unix timestamp in seconds
	EndTimeUnix   int64 // Unix timestamp in seconds, 0 means not ended
	Type          ActivityType
}

// 添加輔助方法來處理時間轉換
func (a *Activity) StartTime() time.Time {
	return time.Unix(a.StartTimeUnix, 0)
}

func (a *Activity) EndTime() time.Time {
	if a.EndTimeUnix == 0 {
		return time.Time{} // return zero time for not ended activities
	}
	return time.Unix(a.EndTimeUnix, 0)
}

func (a *Activity) SetStartTime(t time.Time) {
	a.StartTimeUnix = t.Unix()
}

func (a *Activity) SetEndTime(t time.Time) {
	if t.IsZero() {
		a.EndTimeUnix = 0
	} else {
		a.EndTimeUnix = t.Unix()
	}
}

func (a *Activity) IsEnded() bool {
	return a.EndTimeUnix > 0
}
