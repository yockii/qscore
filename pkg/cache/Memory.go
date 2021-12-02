package cache

import (
	"encoding/json"
	"time"
)

type memoryCacher struct {
	store map[string]*storeValue
}

type storeValue struct {
	expiredAt time.Time
	value     interface{}
}

func (c *memoryCacher) SetWithExpire(key string, value interface{}, expireInSecond int) error {
	if c.store == nil {
		c.store = make(map[string]*storeValue)
	}
	c.store[key] = &storeValue{
		expiredAt: time.Now().Add(time.Second * time.Duration(expireInSecond)),
		value:     value,
	}
	return nil
}

func (c *memoryCacher) GetString(key string) (string, error) {
	if c.store == nil {
		return "", nil
	}
	v := c.store[key]
	if v.expiredAt.Before(time.Now()) {
		delete(c.store, key)
		return "", nil
	}
	s, ok := v.value.(string)
	if ok {
		return s, nil
	}
	bs, err := json.Marshal(v.value)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func (*memoryCacher) Type() string {
	return "memory"
}
func (c *memoryCacher) Close() {
	c.store = make(map[string]*storeValue)
}
