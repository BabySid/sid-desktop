package base

import "time"

func FormatNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func FormatTimeStamp(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02 15:04:05")
}
