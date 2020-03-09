package goviewnet

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/yanzhen74/goview/src/model"
)

type NetRedis struct {
	name        string
	pool        *redis.Pool
	redisServer string
	subscribers map[string]*[]*model.FrameType
}

func newRedisPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
}

func (this *NetRedis) Init(config *model.NetWork) (int, error) {
	fmt.Printf("init redis network")
	return init_redis((*NetRedis)(this), config)
}

func init_redis(this *NetRedis, config *model.NetWork) (int, error) {

	//ips := strings.Split(config.NetWorkIP, ";")
	this.name = config.NetWorkName
	this.redisServer = config.NetWorkIP
	this.pool = newRedisPool(config.NetWorkIP)
	this.subscribers = make(map[string]*[]*model.FrameType)

	return 1, nil
}

func (this *NetRedis) Subscribe(sub *model.FrameType) {
}

// receive net frame, parse and dispatch
func (this *NetRedis) Process() error {
	return nil
}
