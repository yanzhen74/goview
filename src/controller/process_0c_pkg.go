package controller

import (
	"fmt"
	"time"

	"github.com/yanzhen74/goview/src/model"
)

func Process0cPkg(frame model.FrameDict) {
	chan_view := <-frame.Frame_type.ChanViewReg
	frame.Frame_type.ChanViewList = append(frame.Frame_type.ChanViewList, chan_view)
	var pkg string
	fmt.Println("************Got new channel client view first ")
	for i := 0; ; i++ {
		select {
		case chan_view := <-frame.Frame_type.ChanViewReg:
			frame.Frame_type.ChanViewList = append(frame.Frame_type.ChanViewList, chan_view)
			fmt.Println("************Got new channel client view ")
		case <-time.After(time.Millisecond * time.Duration(100)):
			pkg = fmt.Sprintf("0,%d,中国字%d,%d;1,%d,%d,%d;2,%d,%d,%d", i, i, i, i, i, i, i, i, i)
		case frame.Frame_type.ChanViewList[len(frame.Frame_type.ChanViewList)-1] <- pkg:
			pkg = fmt.Sprintf("0,%d,中国字%d,%d;1,%d,%d,%d;2,%d,%d,%d", i, i, i, i, i, i, i, i, i)
			fmt.Printf("channel no %s %d\n", frame.Frame_type.MissionID, len(frame.Frame_type.ChanViewList))
			time.Sleep(time.Millisecond * time.Duration(100))
		}
	}
}
