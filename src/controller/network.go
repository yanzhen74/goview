package controller

import (
	"fmt"

	"github.com/yanzhen74/goview/src/goviewnet"
	"github.com/yanzhen74/goview/src/model"
)

var netProcessers *[]*goviewnet.NetProcesser

func Init_network(conf string) bool {
	// init net config
	netConfig, err := model.Read_network_config(conf)
	if err != nil {
		fmt.Printf("error is %v", err)
		return false
	}
	netProcessers = new([]*goviewnet.NetProcesser)
	for _, n := range (*netConfig).NetWorkList {
		netProcesser := goviewnet.GetNetProcesser(n.NetWorkProtocal)
		ok, _ := (*netProcesser).Init(&n)
		if ok == 1 {
			*netProcessers = append(*netProcessers, netProcesser)
		}
	}
	return true
}

func Bind_network(frame model.FrameType) {
	for _, p := range *netProcessers {
		(*p).Subscribe(&frame)
	}
}

func Run_network() bool {
	for _, p := range *netProcessers {
		go (*p).Process()
	}
	return true
}
