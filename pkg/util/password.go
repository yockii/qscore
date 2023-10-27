package util

import (
	"math/rand"
	"regexp"
)

// PasswordStrengthCheck 校验密码强度，4个级别，0最低，3最高
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
	return l > level
}

// RandomPassword 随机密码，注意此方法不确保字符串唯一性，仅保证根据0-3的级别生成的密码强度
func RandomPassword(min, max, level int) string {
	if min < level {
		min = level
	}
	if max < min {
		max = min + 1
	}

	strSelectList := []string{
		"0123456789",
		"abcdefghijkmnpqrstuvwxyz",
		"ABCDEFGHJKLMNPQRSTUVWXYZ",
		"~!@#$%^&*?_-",
	}
	collectStr := ""

	// 确定每一个级别获取到的字符串长度，确保总长度在min和max之间并且每一个指定级别之下都有字符
	for i := 0; i <= level; i++ {
		// 从strSelectList对应级别中随机抽取一个
		selectFrom := strSelectList[i]
		idx := rand.Intn(len(selectFrom))
		result := selectFrom[idx : idx+1]
		collectStr += result
	}

	left := min + rand.Intn(max-min) - level // 已经有level个字符了，还需要left个字符
	for i := 0; i < left; i++ {
		idx := rand.Intn(level + 1) // 从level个级别中随机抽取一个
		selectFrom := strSelectList[idx]
		idx = rand.Intn(len(selectFrom))
		result := selectFrom[idx : idx+1]
		collectStr += result
	}

	// 得到最终的字符集，进行随机打乱
	result := []byte(collectStr)
	rand.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})
	return string(result)
}
