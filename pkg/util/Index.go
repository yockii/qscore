package util

import (
	"github.com/segmentio/ksuid"
	"strconv"
	"strings"
	"time"
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
