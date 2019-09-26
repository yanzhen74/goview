package controller

import (
	"log"

	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
	"github.com/yanzhen74/goview/src/model"
)

const namespace = "default"

var Frame_page_map map[string][]int
var view_chan_list []chan string
var Dicts *[]model.FrameDict

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
			return nil
		},
		"chat": func(nsConn *websocket.NSConn, msg websocket.Message) error {
			// room.String() returns -> NSConn.String() returns -> Conn.String() returns -> Conn.ID()
			log.Printf("[%s] sent: %s", nsConn, string(msg.Body))
			log.Printf("Server got: %s from [%s]", msg.Body, nsConn.Conn.ID())

			nsConn.Conn.Server().Broadcast(nsConn, msg)

			// add a channel between process_0c_pkg and publishPkg
			view_chan := make(chan string, 10)
			for _, f := range Frame_page_map[(string)(msg.Body)] {
				(*Dicts)[f].Frame_type.ChanViewReg <- view_chan
			}

			log.Printf("Channel bind ok")

			view_chan_list = append(view_chan_list, view_chan)

			go publishPkg(nsConn, msg, view_chan)

			// Write message back to the client message owner with:
			// nsConn.Emit("chat", msg)
			// Write message to all except this client with:
			nsConn.Conn.Server().Broadcast(nsConn, msg)
			return nil
		},
	},
}

func SetupWebsocket(app *iris.Application) {
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

	view_chan_list = make([]chan string, 0, 100)
	// ws.OnDisconnect = func(c *websocket.Conn) {
	// 	log.Printf("[%s] Disconnected from server", c.ID())
	// }
	// register the server on an endpoint.
	// see the inline javascript code in the websockets.html,
	// this endpoint is used to connect to the server.
	app.Get("/echo", websocket.Handler(ws))
	// serve the javascript built'n client-side library,
	// see websockets.html script tags, this path is used.
	// app.Any("/iris-ws.js", websocket.ClientHandler())

}
