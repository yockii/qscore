package domain

import (
	"fmt"
	"reflect"
	"time"
)

var SyncDomains []interface{}

var DateTimeFormat = "2006-01-02 15:04:05"

type DateTime time.Time

func (t DateTime) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("\"%s\"", time.Time(t).Format(DateTimeFormat))
	return []byte(s), nil
}

func (t *DateTime) UnmarshalJSON(data []byte) error {
	s := string(data)
	if s == "null" {
		*t = DateTime(time.Time{})
		return nil
	}
	tm, err := time.Parse(DateTimeFormat, s)
	if err != nil {
		return err
	}
	*t = DateTime(tm)
	return nil
}

var DateTimeConverter = func(value string) reflect.Value {
	if v, err := time.Parse("2006-01-02 15:04:05", value); err == nil {
		return reflect.ValueOf(v)
	}
	return reflect.Value{}
}
