package base

import "time"

const (
	DateTimeFormat = "2006-01-02 15:04:05"
	TimeFormat     = "15:04:05"
	DateFormat     = "2006-01-02"
)

func FormatTime() string {
	return FormatTimeStampWithFormat(time.Now().Unix(), TimeFormat)
}

func FormatDate() string {
	return FormatTimeStampWithFormat(time.Now().Unix(), DateFormat)
}

func FormatTimeStamp(ts int64) string {
	return FormatTimeStampWithFormat(ts, DateTimeFormat)
}

func FormatTimeStampWithFormat(ts int64, format string) string {
	return time.Unix(ts, 0).Format(format)
}
