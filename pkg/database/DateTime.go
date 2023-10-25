package database

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

var DateTimeFormat = "2006-01-02 15:04:05"

type DateTime time.Time

// ////// database/sql 需要的方法
func (t DateTime) String() string {
	return time.Time(t).Format(DateTimeFormat)
}
func (t DateTime) Value() (driver.Value, error) {
	if t.IsZero() {
		return nil, nil
	}
	return time.Time(t), nil
}
func (t *DateTime) Scan(value interface{}) error {
	timeValue, ok := value.(time.Time)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal time value:", value))
	}
	*t = DateTime(timeValue)
	return nil
}

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
