package cache

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

var Prefix string
var Redis *redis.Pool
var enabled bool

func InitRedis(redisPrefix, host, password string, port, maxIdle, maxActive int, options ...redis.DialOption) {
	Prefix = redisPrefix
	Redis = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", host, port), options...)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: 240 * time.Second,
		Wait:        true,
	}
	enabled = true
}

func Get() redis.Conn {
	return Redis.Get()
}

func Close() {
	_ = Redis.Close()
}

func Enabled() bool {
	return enabled
}
