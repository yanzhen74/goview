package model

import "github.com/kataras/iris/websocket"

type View_page_regist_info struct {
	Conn      *websocket.NSConn
	File      string
	View_dict *[]*ViewDict
	View_chan chan string
	Action    int // 1-regist; 0-unregist
}

func Get_view_page_regist_info(view *Paras, view_chan chan string) *View_page_regist_info {
	var info *View_page_regist_info = new(View_page_regist_info)

	info.View_dict = Get_view_page_dict(*view)
	info.File = view.File
	info.View_chan = view_chan
	info.Action = 1

	return info
}

func (info *View_page_regist_info) Set_action(action int) {
	info.Action = action
}
