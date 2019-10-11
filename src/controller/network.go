package controller

import (
	"github.com/yanzhen74/goview/src/model"
	"github.com/yanzhen74/goview/src/net"
)

var netProcessers *[]*net.NetProcesser

func Init_network(config *model.NetWorks) bool {
	netProcessers = new([]*net.NetProcesser)
	for _, n := range (*config).NetWorkList {
		netProcesser := net.GetNetProcesser(n.NetWorkProtocal)
		ok, _ := (*netProcesser).Init(&n)
		if ok == 1 {
			*netProcessers = append(*netProcessers, netProcesser)
		}
	}
	return true
}

func Bind_network(frame model.FrameDict) {
	for _, p := range *netProcessers {
		(*p).Subscribe(frame.Frame_type.DataType, frame.Frame_type.NetChanFrame)
	}
}

func Run_network() bool {
	for _, p := range *netProcessers {
		go (*p).Process()
	}
	return true
}
