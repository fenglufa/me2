package utils

import "time"

const (
	DateFormat     = "2006-01-02"
	TimeFormat     = "15:04:05"
	DateTimeFormat = "2006-01-02 15:04:05"
)

// FormatDate 格式化日期
func FormatDate(t time.Time) string {
	return t.Format(DateFormat)
}

// FormatDateTime 格式化日期时间
func FormatDateTime(t time.Time) string {
	return t.Format(DateTimeFormat)
}

// ParseDate 解析日期
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse(DateFormat, dateStr)
}

// ParseDateTime 解析日期时间
func ParseDateTime(dateTimeStr string) (time.Time, error) {
	return time.Parse(DateTimeFormat, dateTimeStr)
}

// GetToday 获取今天日期字符串
func GetToday() string {
	return time.Now().Format(DateFormat)
}

// GetNow 获取当前时间字符串
func GetNow() string {
	return time.Now().Format(DateTimeFormat)
}

// IsToday 判断是否是今天
func IsToday(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day()
}
