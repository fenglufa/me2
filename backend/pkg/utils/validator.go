package utils

import "regexp"

// IsValidPhone 验证手机号
func IsValidPhone(phone string) bool {
	pattern := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

// IsValidCode 验证验证码
func IsValidCode(code string) bool {
	pattern := `^\d{6}$`
	matched, _ := regexp.MatchString(pattern, code)
	return matched
}

// IsValidEmail 验证邮箱
func IsValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// MaskPhone 手机号脱敏
func MaskPhone(phone string) string {
	if len(phone) != 11 {
		return phone
	}
	return phone[:3] + "****" + phone[7:]
}
