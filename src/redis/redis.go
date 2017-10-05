package redis

import (
	"log"
	"time"

	"github.com/ahmadmuzakki29/redigomock"
	redigo "github.com/garyburd/redigo/redis"
)

type RedisConfig struct {
	EndPoint string
	Timeout  time.Duration
	MaxIdle  int
}

type Redis struct {
	redigo.Pool
	Mock *redigomock.Conn
}

func InitRedis() *Redis {
	var a *RedisConfig = &RedisConfig{
		"devel-redis.tkpd:6379",
		10,
		10,
	}
	var redis *Redis
	if redis == nil {
		newRedis := redigo.Pool{
			MaxIdle:   a.MaxIdle,
			MaxActive: 16,
			Dial: func() (redigo.Conn, error) {
				return redigo.Dial(
					"tcp", a.EndPoint,
					redigo.DialConnectTimeout(a.Timeout*time.Second),
				)
			},
		}

		redis = &Redis{
			newRedis,
			nil,
		}
	}
	return redis
}
func SetRedis(msg interface{}) {
	r := InitRedis()
	result, err := r.Pool.Get().Do("SET", "order_list:Robby", msg)
	if result != "OK" && err != nil {
		log.Println(err)
	}
	log.Println("redis berhasil")
}

func GetRedis(key string) interface{} {
	r := InitRedis()
	result, err := r.Pool.Get().Do("GET", key)
	if err != nil {
		log.Println(err)
	}
	test, _ := redigo.Bytes(result, err)
	return string(test)
}
