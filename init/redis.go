package init

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

func RedisCnnPool(redisHost string, redisPort string, redisPwd string) (pool *redis.Pool) {
	// Redis 连接池
	pool = &redis.Pool{
		MaxIdle:     100,
		MaxActive:   4000,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisHost+":"+redisPort, redis.DialPassword(redisPwd))
			if nil != err {
				StdPrint.Error("redis connect error: ", err)
				return nil, err
			}
			StdPrint.Info("redis connected")
			return c, nil
		},
	}
	return
}

func RedisIsExist(key string) (result bool) {
	result = false
	cnn := RedisClient.Get()
	defer cnn.Close()

	result, _ = redis.Bool(cnn.Do("EXISTS", key))

	return result
}

func RedisSetString(key string, value string, seconds int) (e error) {
	cnn := RedisClient.Get()
	defer cnn.Close()

	_, e = cnn.Do("SET", key, value, "EX", seconds)
	return
}

func RedisGetString(key string) (value string, e error) {
	cnn := RedisClient.Get()
	defer cnn.Close()

	value, e = redis.String(cnn.Do("GET", key))
	return
}

func RedisDelString(key string) (e error) {
	cnn := RedisClient.Get()
	defer cnn.Close()

	_, e = redis.String(cnn.Do("DEL", key))
	return
}

func RedisTTL(key string) (ttl int64, e error) {
	cnn := RedisClient.Get()
	defer cnn.Close()

	ttl, e = redis.Int64(cnn.Do("TTL", key))
	return
}
