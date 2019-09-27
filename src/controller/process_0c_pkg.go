package controller

import (
	"fmt"
	"reflect"
	"time"

	"github.com/yanzhen74/goview/src/model"
)

func Process0cPkg(frame model.FrameDict) {
	// chan_view := <-frame.Frame_type.ChanViewReg
	// frame.Frame_type.ChanViewList = append(frame.Frame_type.ChanViewList, chan_view)
	var pkg string
	fmt.Println("************Got new channel client view first ", pkg)
	ticker := time.NewTicker(time.Millisecond * time.Duration(100))
	cases := update(ticker, frame.Frame_type.ChanViewReg, frame.Frame_type.ChanViewList, pkg)

	for i := 0; ; i++ {
		chose, value, _ := reflect.Select(cases)
		switch chose {
		case 0: // new chan_view
			info := (value.Interface().(*model.View_page_regist_info))
			frame.Frame_type.ChanViewList = append(frame.Frame_type.ChanViewList, info.View_chan)
			cases = update(ticker, frame.Frame_type.ChanViewReg, frame.Frame_type.ChanViewList, pkg)
		case 1: // time
			pkg = fmt.Sprintf("0,%d,中国字%d,%d;1,%d,%d,%d;2,%d,%d,%d", i, i, i, i, i, i, i, i, i)
		default:
			pkg = fmt.Sprintf("0,%d,%s%d,%d;1,%d,%d,%d;2,%d,%d,%d", i, frame.Frame_type.MissionID, i, i, i, i, i, i, i, i)
			// fmt.Printf("channel no %s %d\n", frame.Frame_type.MissionID, len(frame.Frame_type.ChanViewList))
			time.Sleep(time.Millisecond * time.Duration(100))
			cases = update(ticker, frame.Frame_type.ChanViewReg, frame.Frame_type.ChanViewList, pkg)
		}
	}
	// for i := 0; ; i++ {
	// 	select {
	// 	case chan_view := <-frame.Frame_type.ChanViewReg:
	// 		frame.Frame_type.ChanViewList = append(frame.Frame_type.ChanViewList, chan_view)
	// 		fmt.Println("************Got new channel client view ")
	// 	case <-time.After(time.Millisecond * time.Duration(100)):
	// 		pkg = fmt.Sprintf("0,%d,中国字%d,%d;1,%d,%d,%d;2,%d,%d,%d", i, i, i, i, i, i, i, i, i)
	// 	case frame.Frame_type.ChanViewList[len(frame.Frame_type.ChanViewList)-1] <- pkg:
	// 		pkg = fmt.Sprintf("0,%d,中国字%d,%d;1,%d,%d,%d;2,%d,%d,%d", i, i, i, i, i, i, i, i, i)
	// 		fmt.Printf("channel no %s %d\n", frame.Frame_type.MissionID, len(frame.Frame_type.ChanViewList))
	// 		time.Sleep(time.Millisecond * time.Duration(100))
	// 	}
	// }
}

func update(
	ticker *time.Ticker,
	chan_view_reg chan *model.View_page_regist_info,
	chan_view_list []chan string,
	send_value interface{}) (cases []reflect.SelectCase) {
	// chan view register
	selectcase := reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(chan_view_reg),
	}
	cases = append(cases, selectcase)

	// 定时器
	selectcase = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(ticker.C),
	}
	cases = append(cases, selectcase)

	// 每个消费者
	for _, item := range chan_view_list {
		selectcase = reflect.SelectCase{
			Dir:  reflect.SelectSend,
			Chan: reflect.ValueOf(item),
			Send: reflect.ValueOf(send_value),
		}
		cases = append(cases, selectcase)
	}
	return
}
