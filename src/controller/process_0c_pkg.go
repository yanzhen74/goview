package controller

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/kataras/iris/websocket"
	"github.com/yanzhen74/goview/src/model"
)

// map from para index of frame to para index of view
var para_view_map map[*websocket.NSConn]map[int]string

func Process0cPkg(frame model.FrameDict) {
	// chan_view := <-frame.Frame_type.UserChanReg
	// frame.Frame_type.UserChanMap = append(frame.Frame_type.UserChanMap, chan_view)
	para_view_map = make(map[*websocket.NSConn]map[int]string)

	// pkg should send only required parameters to view's chan
	var pkg map[chan string]interface{} = make(map[chan string]interface{})

	ticker := time.NewTicker(time.Millisecond * time.Duration(100))
	cases := update(ticker, frame.Frame_type.UserChanReg, frame.Frame_type.UserChanMap, pkg)

	for i := 0; ; {
		chose, value, _ := reflect.Select(cases)

		// update when receive net data
		v := fmt.Sprintf("0,%d,%s%d,%d;1,%d,%d,%d;2,%d,%d,%d", i, frame.Frame_type.MissionID, i, i, i, i, i, i, i, i)
		for _, c := range frame.Frame_type.UserChanMap {
			pkg[c] = v
		}

		// log.Printf("cases len %d, channel no %s %d %d\n", len(cases), frame.Frame_type.MissionID, len(frame.Frame_type.UserChanMap), i)
		switch chose {
		case 0: // regist/unregist chan_view
			info := (value.Interface().(*model.View_page_regist_info))
			if regist_view_chan(&frame, info) {
				pkg[info.View_chan] = ""
				cases = update(ticker, frame.Frame_type.UserChanReg, frame.Frame_type.UserChanMap, pkg)
			}
		case 1: // time
			cases = update(ticker, frame.Frame_type.UserChanReg, frame.Frame_type.UserChanMap, pkg)
			i++
		default:
			// pkg = fmt.Sprintf("0,%d,%s%d,%d;1,%d,%d,%d;2,%d,%d,%d", i, frame.Frame_type.MissionID, i, i, i, i, i, i, i, i)
			//fmt.Printf("send ok %s %d %d\n", frame.Frame_type.MissionID, len(frame.Frame_type.UserChanMap), i)
			// fmt.Printf("channel no %s %d %d\n", frame.Frame_type.MissionID, len(frame.Frame_type.UserChanMap), i)
			time.Sleep(time.Millisecond * time.Duration(10))
			//cases = update(ticker, frame.Frame_type.UserChanReg, frame.Frame_type.UserChanMap, pkg)
		}
	}
}

func regist_view_chan(frame *model.FrameDict, info *model.View_page_regist_info) bool {
	if info.Action == 1 {
		// regist only required parameters for view
		para_view_map[info.Conn] = make(map[int]string)

		var view_dict *model.ViewDict = nil
		for _, v := range *(info.View_dict) {
			if (*v).View_type.PayloadName == frame.Frame_type.PayloadName && (*v).View_type.SubAddressName == frame.Frame_type.SubAddressName {
				view_dict = v
				break
			}
		}
		if view_dict == nil {
			log.Printf("Failed register New channel no %s %s %s, now %d user\n", frame.Frame_type.MissionID, info.Conn, info.File,
				len(frame.Frame_type.UserChanMap))
			return false
		}
		for index, item := range frame.ParaList {
			for _, p := range (*view_dict).ParaList {
				if p.ParaKey == item.ParaKey {
					(para_view_map[info.Conn])[index] = p.Index
					log.Printf("bound %s\n", p.Index)
				}
			}
		}

		// bind from conn to view_chan
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
	return true
}

func update(
	ticker *time.Ticker,
	chan_view_reg chan *model.View_page_regist_info,
	user_chan_view_map map[*websocket.NSConn]chan string,
	send_value_map map[chan string]interface{}) (cases []reflect.SelectCase) {

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
		send_value := send_value_map[item]
		selectcase = reflect.SelectCase{
			Dir:  reflect.SelectSend,
			Chan: reflect.ValueOf(item),
			Send: reflect.ValueOf(send_value),
		}
		cases = append(cases, selectcase)
	}
	return
}
