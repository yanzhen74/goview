package goviewnet

import (
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	pool        *redis.Pool
	redisServer = flag.String("10.211.55.4", ":6379", "")
)

func Test_connect(t *testing.T) {
	c, err := redis.Dial("tcp", "10.211.55.4:6379")
	if err != nil {
		fmt.Println("connect to redis error", err)
		return
	}
	defer c.Close()
}

func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
}

func Test_redis_pool(t *testing.T) {
	flag.Parse()
	pool = newPool("10.211.55.4:6379")

	c := pool.Get()
	defer c.Close()

	_, err := c.Do("SET", "mykey", "what a pool")
	if err != nil {
		fmt.Println("set to redis error", err)
		return
	}

	value, err := redis.String(c.Do("GET", "mykey"))
	if err != nil {
		fmt.Println("get from redis error", err)
	} else {
		fmt.Printf("Get from redis %v", value)
	}

	pool.Close()
}

func Test_redis_consumer(t *testing.T) {
	c, err := redis.Dial("tcp", "10.211.55.4:6379")
	if err != nil {
		fmt.Println("connect to redis error", err)
		return
	}
	defer c.Close()

	_, err = c.Do("SET", "mykey", "what a fuck")
	if err != nil {
		fmt.Println("set to redis error", err)
		return
	}

	value, err := redis.String(c.Do("GET", "mykey"))
	if err != nil {
		fmt.Println("get from redis error", err)
	} else {
		fmt.Printf("Get from redis %v", value)
	}
}
func Test_redis_list(t *testing.T) {
	c, err := redis.Dial("tcp", "10.211.55.4:6379")
	if err != nil {
		fmt.Println("connect to redis error", err)
		return
	}
	defer c.Close()

	_, err = c.Do("lpush", "dbs", "redis")
	if err != nil {
		fmt.Println("set to redis error", err)
		return
	}

	_, err = c.Do("lpush", "dbs", "mysql")
	if err != nil {
		fmt.Println("set to redis error", err)
		return
	}
	_, err = c.Do("lpush", "dbs", "sqlite")
	if err != nil {
		fmt.Println("set to redis error", err)
		return
	}

	values, _ := redis.Values(c.Do("lrange", "dbs", "0", "20"))

	for _, v := range values {
		fmt.Println(string(v.([]byte)))
	}

}
