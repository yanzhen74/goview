package routes

import (
	"log"

	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
	"github.com/yanzhen74/goview/src/controller"
	"github.com/yanzhen74/goview/src/model"
)

const namespace = "default"

// map nsConn to info
var conn_info_map map[*websocket.NSConn]*model.View_page_regist_info

func regist_info(nsConn *websocket.NSConn, action int) {
	info := conn_info_map[nsConn]
	info.Set_action(action)
	for _, i := range *(info.View_dict) {
		for _, d := range *controller.Dicts {
			if d.Frame_type.MissionID == (*i).View_type.MissionID &&
				d.Frame_type.DataType == (*i).View_type.DataType &&
				d.Frame_type.PayloadName == (*i).View_type.PayloadName &&
				d.Frame_type.SubAddressName == (*i).View_type.SubAddressName {
				d.Frame_type.UserChanReg <- info
				if action == 0 {
					log.Printf("Channel unbound %s ok\n", info.File)
				} else {
					log.Printf("Channel bound %s ok\n", info.File)
				}
				break
			}
		}
	}
}

// if namespace is empty then simply websocket.Events{...} can be used instead.
var serverEvents = websocket.Namespaces{
	namespace: websocket.Events{
		websocket.OnNamespaceConnected: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			// with `websocket.GetContext` you can retrieve the Iris' `Context`.
			ctx := websocket.GetContext(nsConn.Conn)

			log.Printf("[%s] connected to namespace [%s] with IP [%s]",
				nsConn, msg.Namespace,
				ctx.RemoteAddr())
			return nil
		},
		websocket.OnNamespaceDisconnect: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			log.Printf("[%s] disconnected from namespace [%s]", nsConn, msg.Namespace)
			regist_info(nsConn, 0)

			return nil
		},
		"chat": func(nsConn *websocket.NSConn, msg websocket.Message) error {
			// room.String() returns -> NSConn.String() returns -> Conn.String() returns -> Conn.ID()
			log.Printf("[%s] sent: %s", nsConn, string(msg.Body))
			log.Printf("Server got: %s from [%s]", msg.Body, nsConn.Conn.ID())

			nsConn.Conn.Server().Broadcast(nsConn, msg)

			// bind frontend view page to backend frame which contains their params
			paras := File_paras_map[(string)(msg.Body)]
			// add a channel between backend and frontend
			view_chan := make(chan string, 10)
			info := model.Get_view_page_regist_info(paras, view_chan)

			// bind nsConn to info, for unregist it when close
			conn_info_map[nsConn] = info
			info.Conn = nsConn

			// regist info
			regist_info(nsConn, 1)

			go controller.PublishPkg(nsConn, msg, view_chan)

			// Write message back to the client message owner with:
			// nsConn.Emit("chat", msg)
			// Write message to all except this client with:
			// nsConn.Conn.Server().Broadcast(nsConn, msg)
			return nil
		},
	},
}

func WebsocketHub(party iris.Party) {
	web := party.Party("/echo")

	// create our websocket server
	// Almost all features of neffos are disabled because no custom message can pass
	// when app expects to accept and send only raw websocket native messages.
	// When only allow native messages is a fact?
	// When the registered namespace is just one and it's empty
	// and contains only one registered event which is the `OnNativeMessage`.
	// When `Events{...}` is used instead of `Namespaces{ "namespaceName": Events{...}}`
	// then the namespace is empty "".
	ws := websocket.New(
		websocket.DefaultGorillaUpgrader,
		serverEvents)
	// ws.OnConnect = func(c *websocket.Conn) error {
	// 	log.Printf("[%s] Connected to server!", c.ID())
	// 	return nil
	// }

	conn_info_map = make(map[*websocket.NSConn]*model.View_page_regist_info)
	// ws.OnDisconnect = func(c *websocket.Conn) {
	// 	log.Printf("[%s] Disconnected from server", c.ID())
	// }
	// register the server on an endpoint.
	// see the inline javascript code in the websockets.html,
	// this endpoint is used to connect to the server.
	web.Get("/", websocket.Handler(ws)) // websocket模块
	//app.Get("/echo", websocket.Handler(ws))
	// serve the javascript built'n client-side library,
	// see websockets.html script tags, this path is used.
	// app.Any("/iris-ws.js", websocket.ClientHandler())

}
