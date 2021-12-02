package cache

import (
	"github.com/gomodule/redigo/redis"

	"github.com/yockii/qscore/pkg/logger"
)

type redisCacher struct {
	redisPrefix string
	redisPool   *redis.Pool
}

func (c *redisCacher) SetWithExpire(key string, value interface{}, expireInSecond int) error {
	redisConn := c.redisPool.Get()
	defer redisConn.Close()
	_, err := redisConn.Do("SETEX", key, expireInSecond, value)
	if err != redis.ErrNil {
		return err
	}
	return nil
}

func (c *redisCacher) GetString(key string) (string, error) {
	redisConn := c.redisPool.Get()
	defer redisConn.Close()
	v, err := redis.String(redisConn.Do("GET", key))
	if err != nil {
		return "", err
	}
	return v, nil
}

func (*redisCacher) Type() string {
	return "redis"
}
func (c *redisCacher) Close() {
	if err := c.redisPool.Close(); err != nil {
		logger.Error(err)
	}
}
