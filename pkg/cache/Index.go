package cache

import (
	"errors"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

type cacher interface {
	Type() string
	Close()
	GetString(string) (string, error)
}

type cacheable struct {
	inMemory bool
	mc       cacher
	// redis
	useRedis bool
	rc       cacher
}

var cacheInstance = &cacheable{}

func UseMemory() {
	cacheInstance.inMemory = true
	var c cacher = &memoryCacher{}
	cacheInstance.mc = c
}

func UseRedis(redisPrefix, host, password string, port, maxIdle, maxActive int) {
	var rc cacher = &redisCacher{
		redisPrefix: redisPrefix,
		redisPool: &redis.Pool{
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
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
		},
	}
	cacheInstance.rc = rc
}

func Close() {
	if cacheInstance.inMemory {
		cacheInstance.mc.Close()
	}
	if cacheInstance.useRedis {
		cacheInstance.rc.Close()
	}
}

func GetString(key string) (value string, err error) {
	if cacheInstance.useRedis {
		return cacheInstance.rc.GetString(key)
	}
	if cacheInstance.inMemory {
		return cacheInstance.mc.GetString(key)
	}
	return "", errors.New("GetString方法未找到")
}
