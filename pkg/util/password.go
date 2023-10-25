package util

import "regexp"

// PasswordStrengthCheck 校验密码强度，5个级别，0最低，4最高
func PasswordStrengthCheck(min, max, level int, pwd string) bool {
	if len(pwd) < min {
		return false
	}
	if len(pwd) > max {
		return false
	}
	l := 0
	patternList := []string{`[0-9]+`, `[a-z]+`, `[A-Z]+`, `[~!@#$%^&*?_-]+`}
	for _, pattern := range patternList {
		match, _ := regexp.MatchString(pattern, pwd)
		if match {
			l++
		}
	}
	return l >= level
}
