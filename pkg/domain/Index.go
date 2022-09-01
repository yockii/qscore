package domain

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

var SyncDomains []interface{}

var DateTimeFormat = "2006-01-02 15:04:05"

type DateTime time.Time

func (t DateTime) IsZero() bool {
	return time.Time(t).IsZero()
}

func (t DateTime) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	s := fmt.Sprintf("\"%s\"", time.Time(t).Format(DateTimeFormat))
	return []byte(s), nil
}

func (t *DateTime) UnmarshalJSON(data []byte) error {
	s := string(data)
	s = strings.Replace(s, "\"", "", -1)
	if s == "null" || s == "" {
		*t = DateTime(time.Time{})
		return nil
	}
	tm, err := time.ParseInLocation(DateTimeFormat, s, time.Local)
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
	return reflect.ValueOf(time.Time{})
}

type TimeCondition struct {
	Start DateTime `json:"start,omitempty" query:"start"`
	End   DateTime `json:"end,omitempty" query:"end"`
}

type CommonResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

type Paginate struct {
	Total  int         `json:"total"`
	Offset int         `json:"offset"`
	Limit  int         `json:"limit"`
	Items  interface{} `json:"items"`
}
