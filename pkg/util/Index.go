package util

import (
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/ksuid"
)

var CstZone = time.FixedZone("CST", 8*3600) // 东八

func GenerateDatabaseID() string {
	uid := ksuid.New()
	return uid.String()
}

func Unicode2Zh(form string) (to string, err error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(form), `\\u`, `\u`, -1))
	if err != nil {
		return "", err
	}
	return str, nil
}

// SnakeString XxYy to xx_yy , XxYY to xx_yy
func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

// CamelString xx_yy to XxYy
func CamelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}
