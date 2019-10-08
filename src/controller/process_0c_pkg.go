package controller

import (
	"bytes"
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
	cases := update(frame.Frame_type.NetChanFrame, ticker, frame.Frame_type.UserChanReg, nil, nil)

	d := 0
	for i := 0; ; {
		chose, value, _ := reflect.Select(cases)

		// log.Printf("cases len %d, channel no %s %d %d\n", len(cases), frame.Frame_type.MissionID, len(frame.Frame_type.UserChanMap), i)
		switch chose {
		case 0: // regist/unregist chan_view
			info := (value.Interface().(*model.View_page_regist_info))
			if 1 == regist_view_chan(&frame, info) {
				pkg[info.View_chan] = ""
			}
			cases = update(frame.Frame_type.NetChanFrame, ticker, frame.Frame_type.UserChanReg, nil, nil)
		case 1: // time
			frame.Frame_type.NetChanFrame <- "hello world"
		case 2: // net frame
			// update when receive net data
			v := make([]string, 0, 10)
			for j := 0; j < len(frame.ParaList); j++ {
				v = append(v, fmt.Sprintf(",%d,%s%d,%d;", i, frame.Frame_type.MissionID, i, i))
			}

			var buffer bytes.Buffer
			for conn, view_chan := range frame.Frame_type.UserChanMap {
				buffer.Reset()
				for id_in_frame, id_in_view := range para_view_map[conn] {
					buffer.WriteString(id_in_view)
					buffer.WriteString(v[id_in_frame])
				}
				pkg[view_chan] = buffer.String()
			}
			cases = update(frame.Frame_type.NetChanFrame, ticker, frame.Frame_type.UserChanReg, frame.Frame_type.UserChanMap, pkg)
			i++
		default:
			// log.Printf("default %d\n", d)
			cases = append(cases[:chose], cases[chose+1:]...)
			d++
			// pkg = fmt.Sprintf("0,%d,%s%d,%d;1,%d,%d,%d;2,%d,%d,%d", i, frame.Frame_type.MissionID, i, i, i, i, i, i, i, i)
			//fmt.Printf("send ok %s %d %d\n", frame.Frame_type.MissionID, len(frame.Frame_type.UserChanMap), i)
			// fmt.Printf("channel no %s %d %d\n", frame.Frame_type.MissionID, len(frame.Frame_type.UserChanMap), i)
			// time.Sleep(time.Millisecond * time.Duration(1))
			//cases = update(ticker, frame.Frame_type.UserChanReg, frame.Frame_type.UserChanMap, pkg)
		}
	}
}

// return 0: not changed; 1: new regist; -1: unregist
func regist_view_chan(frame *model.FrameDict, info *model.View_page_regist_info) int {
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
			return 0
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
		return 1
	} else {
		// unregist para for view, sub_map will be GC autoly
		delete(para_view_map, info.Conn)

		// unbind conn with view_chan
		delete(frame.Frame_type.UserChanMap, info.Conn)
		// 关闭chan
		close(info.View_chan)
		log.Printf("Delete channel no %s %s %s, now %d user\n", frame.Frame_type.MissionID, info.Conn, info.File,
			len(frame.Frame_type.UserChanMap))
		return -1
	}
}

func update(
	chan_net_frame chan string,
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

	// chan view register
	selectcase = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(chan_net_frame),
	}
	cases = append(cases, selectcase)

	// 每个消费者
	if user_chan_view_map == nil {
		return
	}

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
