package util

import (
	"strconv"
	"strings"
	"time"

	"github.com/rs/xid"
	"github.com/segmentio/ksuid"
)

func GenerateDatabaseID() string {
	uid := ksuid.New()
	return uid.String()
}

func GetTimeFromDatabaseId(id string) (t time.Time, err error) {
	kid, err := ksuid.Parse(id)
	if err != nil {
		return
	}
	t = kid.Time()
	return
}

func GenerateRequestID() string {
	uid := xid.New()
	return uid.String()
}

func GetTimeFromRequestId(id string) (t time.Time, err error) {
	x, err := xid.FromString(id)
	if err != nil {
		return
	}
	t = x.Time()
	return
}

func Unicode2Zh(form string) (to string, err error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(form), `\\u`, `\u`, -1))
	if err != nil {
		return "", err
	}
	return str, nil
}
