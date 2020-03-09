package goviewnet

import (
	"fmt"

	"github.com/yanzhen74/goview/src/model"
)

type NetKafkaENG NetKafka

func (this *NetKafkaENG) Init(config *model.NetWork) (int, error) {
	fmt.Printf("init kafka-eng network")
	return init_kafka((*NetKafka)(this), config)
}

func (this *NetKafkaENG) Subscribe(sub *model.FrameType) {
}

// receive net frame, parse and dispatch
func (this *NetKafkaENG) Process() error {
	return nil
}
