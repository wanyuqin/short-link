package timex

import "time"

func FormatDateTime(t int64) string {
	return time.UnixMilli(t).Format("2006-01-02 15:04:05")
}
