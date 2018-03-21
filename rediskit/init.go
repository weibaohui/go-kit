package rediskit

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/go-redis/redis"
	"github.com/weibaohui/go-kit/propkit"
)

var globalRedis *redis.Client
var once sync.Once
var globalRedisKit *redisKit

func getRedis() *redis.Client {
	host := propkit.Init().Get("redis.host")
	port := propkit.Init().GetInt("redis.port")
	password := propkit.Init().Get("redis.password")
	db := propkit.Init().GetInt("redis.db")
	once.Do(func() {
		globalRedis = redis.NewClient(&redis.Options{
			Addr:     host + ":" + strconv.Itoa(port),
			Password: password,
			DB:       db,
		})

		pong, err := globalRedis.Ping().Result()
		fmt.Println(pong, err)
		// Output: PONG <nil>
	})
	return globalRedis
}

type redisKit struct {
	DB *redis.Client
}

func Init() *redisKit {
	if globalRedis == nil {
		globalRedisKit = &redisKit{
			DB: getRedis(),
		}
	}
	return globalRedisKit
}

func InitWith(ip string, port int, password string) *redisKit {

	once.Do(func() {
		globalRedis = redis.NewClient(&redis.Options{
			Addr:     ip + ":" + strconv.Itoa(port),
			Password: password, // no password set
			DB:       0,        // use default DB
		})
		pong, err := globalRedis.Ping().Result()
		fmt.Println(pong, err)
		// Output: PONG <nil>
	})
	redisKit := redisKit{
		DB: getRedis(),
	}
	return &redisKit
}
