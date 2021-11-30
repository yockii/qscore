package cache

import "github.com/gomodule/redigo/redis"

type redisCacher struct {
	redisPrefix string
	redisPool   *redis.Pool
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
