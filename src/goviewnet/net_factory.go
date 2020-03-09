package goviewnet

import (
	"github.com/yanzhen74/goview/src/model"
)

type NetProcesser interface {
	Init(*model.NetWork) (int, error)
	Process() error
	Subscribe(*model.FrameType)
}

func GetNetProcesser(netType string) *NetProcesser {
	var processer NetProcesser
	switch netType {
	case "kafka":
		processer = new(NetKafka)
	case "kafka_gwg":
		processer = new(NetKafkaGWG)
	case "kafka_eng":
		processer = new(NetKafkaENG)
	default:
	}

	return &processer
}
