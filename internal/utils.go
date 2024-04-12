package internal

import (
	"fmt"
	"time"
)

func FormatDate(date time.Time) string {
	y, m, d := date.Date()
	hours := date.Hour()
	minutes := date.Minute()

	return fmt.Sprintf("%d-%02d-%02d %02d:%02d ", y, m, d, hours, minutes)
}
