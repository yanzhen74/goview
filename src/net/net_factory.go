package net

import (
	"github.com/yanzhen74/goview/src/model"
)

type NetProcesser interface {
	Init(*model.NetWork) (int, error)
	Process() error
	Subscribe(string, chan string)
}

func GetNetProcesser(netType string) *NetProcesser {
	var processer NetProcesser
	switch netType {
	case "kafka":
		processer = new(NetKafka)
	default:
	}

	return &processer
}
