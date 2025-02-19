package utils

import (
	"fmt"
	"time"
)

func FormatDuration(d time.Duration) string {
	h := d.Hours()
	m := d.Minutes() - float64(int(h))*60
	s := d.Seconds() - float64(int(d.Minutes()))*60
	return fmt.Sprintf("%02d:%02d:%02d", int(h), int(m), int(s))
} 