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

	if seconds == -1 {
		// set seconds = -1 for permanent storage
		_, e = cnn.Do("SET", key, value)
	} else {
		_, e = cnn.Do("SET", key, value, "EX", seconds)
	}
	return
}

func RedisGetString(key string) (value string, e error) {
	cnn := RedisClient.Get()
	defer cnn.Close()

	value, e = redis.String(cnn.Do("GET", key))
	return
}

func RedisSetBytes(key string, value []byte, seconds int) (e error) {
	cnn := RedisClient.Get()
	defer cnn.Close()

	if seconds == -1 {
		// set seconds = -1 for permanent storage
		_, e = cnn.Do("SET", key, value)
	} else {
		_, e = cnn.Do("SET", key, value, "EX", seconds)
	}
	return
}

func RedisGetBytes(key string) (value []byte, e error) {
	cnn := RedisClient.Get()
	defer cnn.Close()

	value, e = redis.Bytes(cnn.Do("GET", key))
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

func RedisIsMember(key string, value string) (result int, e error) {
	cnn := RedisClient.Get()
	defer cnn.Close()

	result, e = redis.Int(cnn.Do("SISMEMBER", key, value))
	return
}

func RedisIncrease(key string) (newValue int64, e error) {
	cnn := RedisClient.Get()
	defer cnn.Close()

	newValue, e = redis.Int64(cnn.Do("INCR", key))
	return
}

func RedisZAdd(key string, score string, name string) (added int, e error) {
	cnn := RedisClient.Get()
	defer cnn.Close()

	added, e = redis.Int(cnn.Do("zadd", key, score, name))
	return
}

func RedisZIncrBy(key string, incr string, name string) (newValue int, e error) {
	cnn := RedisClient.Get()
	defer cnn.Close()

	newValue, e = redis.Int(cnn.Do("zincrby", key, incr, name))
	return
}

func RedisZScore(key string, name string) (score string, e error) {
	cnn := RedisClient.Get()
	defer cnn.Close()

	score, e = redis.String(cnn.Do("zscore", key, name))
	return
}

func RedisZCount(key string, min string, max string) (count int, e error) {
	cnn := RedisClient.Get()
	defer cnn.Close()

	count, e = redis.Int(cnn.Do("zcount", key, min, max))
	return
}

func RedisZRange(key string, start string, stop string) (res [][]byte, e error) {
	cnn := RedisClient.Get()
	defer cnn.Close()

	res, e = redis.ByteSlices(cnn.Do("zrange", key, start, stop, "WITHSCORES"))
	return
}

func RedisZRevRange(key string, start string, stop string) (res [][]byte, e error) {
	cnn := RedisClient.Get()
	defer cnn.Close()

	res, e = redis.ByteSlices(cnn.Do("zrevrange", key, start, stop, "WITHSCORES"))
	return
}

func RedisByteArrayToDataset(listBytes [][]byte, err error) (dataset []interface{}, e error) {
	if err != nil {
		e = err
		return
	}

	for i := 0; i < len(listBytes); i = i + 2 {
		value, _ := strconv.ParseInt(string(listBytes[i+1]), 10, 64)
		var newData = struct {
			Name  string `json:"name" q:",名称"`
			Value int64  `json:"value" q:",值"`
		}{
			Name:  string(listBytes[i]),
			Value: value,
		}
		dataset = append(dataset, newData)
	}
	return
}

func RedisByteArrayToDataset(listBytes [][]byte, err error) (dataset []interface{}, e error) {
	if err != nil {
		e = err
		return
	}

	for i := 0; i < len(listBytes); i = i + 2 {
		value, _ := strconv.ParseInt(string(listBytes[i+1]), 10, 64)
		var newData = struct {
			Name  string `json:"name" q:",名称"`
			Value int64  `json:"value" q:",值"`
		}{
			Name:  string(listBytes[i]),
			Value: value,
		}
		dataset = append(dataset, newData)
	}
	return
}


func RedisZRank(key string, name string) (rank int, e error) {
	cnn := RedisClient.Get()
	defer cnn.Close()

	rank, e = redis.Int(cnn.Do("zrank", key, name))
	return
}

func RedisZRevRank(key string, name string) (revRank int, e error) {
	cnn := RedisClient.Get()
	defer cnn.Close()

	revRank, e = redis.Int(cnn.Do("zrevrank", key, name))
	return
}
