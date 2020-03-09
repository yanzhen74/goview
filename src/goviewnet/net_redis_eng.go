package goviewnet

import (
	"fmt"
	"log"

	"github.com/yanzhen74/goview/src/model"
)

type NetRedisENG NetRedis

func (this *NetRedisENG) Init(config *model.NetWork) (int, error) {
	fmt.Printf("init redis-eng network")
	return init_redis((*NetRedis)(this), config)
}

func (this *NetRedisENG) Subscribe(sub *model.FrameType) {
	if this.name == sub.DataType {
		frame_type := get_frame_type(sub)
		log.Println(frame_type)
		if _, ok := this.subscribers[frame_type]; !ok {
			this.subscribers[frame_type] = new([]*model.FrameType)
		}
		*(this.subscribers[frame_type]) = append(*(this.subscribers[frame_type]), sub)
	}
}

// receive net frame, parse and dispatch
func (this *NetRedisENG) Process() error {
	return nil
}
