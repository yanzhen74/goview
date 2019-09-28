package controller

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/kataras/iris/websocket"
	"github.com/yanzhen74/goview/src/model"
)

func Process0cPkg(frame model.FrameDict) {
	// chan_view := <-frame.Frame_type.UserChanReg
	// frame.Frame_type.UserChanMap = append(frame.Frame_type.UserChanMap, chan_view)
	var pkg string
	ticker := time.NewTicker(time.Millisecond * time.Duration(100))
	cases := update(ticker, frame.Frame_type.UserChanReg, frame.Frame_type.UserChanMap, pkg)

	for i := 0; ; {
		chose, value, _ := reflect.Select(cases)
		switch chose {
		case 0: // regist/unregist chan_view
			info := (value.Interface().(*model.View_page_regist_info))
			regist_view_chan(&frame, info)
			cases = update(ticker, frame.Frame_type.UserChanReg, frame.Frame_type.UserChanMap, pkg)
		case 1: // time
			pkg = fmt.Sprintf("0,%d,%s%d,%d;1,%d,%d,%d;2,%d,%d,%d", i, frame.Frame_type.MissionID, i, i, i, i, i, i, i, i)
			// log.Printf("cases len %d, channel no %s %d %d\n", len(cases), frame.Frame_type.MissionID, len(frame.Frame_type.UserChanMap), i)
			cases = update(ticker, frame.Frame_type.UserChanReg, frame.Frame_type.UserChanMap, pkg)
			i++
		default:
			pkg = fmt.Sprintf("0,%d,%s%d,%d;1,%d,%d,%d;2,%d,%d,%d", i, frame.Frame_type.MissionID, i, i, i, i, i, i, i, i)
			//fmt.Printf("send ok %s %d %d\n", frame.Frame_type.MissionID, len(frame.Frame_type.UserChanMap), i)
			// fmt.Printf("channel no %s %d %d\n", frame.Frame_type.MissionID, len(frame.Frame_type.UserChanMap), i)
			time.Sleep(time.Millisecond * time.Duration(10))
			//cases = update(ticker, frame.Frame_type.UserChanReg, frame.Frame_type.UserChanMap, pkg)
		}
	}
}

func regist_view_chan(frame *model.FrameDict, info *model.View_page_regist_info) {
	if info.Action == 1 {
		frame.Frame_type.UserChanMap[info.Conn] = info.View_chan
		log.Printf("New channel no %s %s %s, now %d user\n", frame.Frame_type.MissionID, info.Conn, info.File,
			len(frame.Frame_type.UserChanMap))
	} else {
		delete(frame.Frame_type.UserChanMap, info.Conn)
		// 关闭chan
		close(info.View_chan)
		log.Printf("Delete channel no %s %s %s, now %d user\n", frame.Frame_type.MissionID, info.Conn, info.File,
			len(frame.Frame_type.UserChanMap))
	}
}

func update(
	ticker *time.Ticker,
	chan_view_reg chan *model.View_page_regist_info,
	user_chan_view_map map[*websocket.NSConn]chan string,
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
	for _, item := range user_chan_view_map {
		selectcase = reflect.SelectCase{
			Dir:  reflect.SelectSend,
			Chan: reflect.ValueOf(item),
			Send: reflect.ValueOf(send_value),
		}
		cases = append(cases, selectcase)
	}
	return
}
