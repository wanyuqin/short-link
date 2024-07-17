package timex

import "time"

func FormatDateTime(t int64) string {
	return time.UnixMilli(t).Format("2006-01-02 15:04:05")
}

// GetDayOfStartUnix 获取传入时间的当天开始的时间戳
func GetDayOfStartUnix(t time.Time) int64 {
	startOfDay := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return startOfDay.Unix()
}

func GetDayOfStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

}
