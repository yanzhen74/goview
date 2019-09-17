package main

import (
	"log"

	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
)

func setupWebsocket(app *iris.Application) {
	// create our websocket server
	// Almost all features of neffos are disabled because no custom message can pass
	// when app expects to accept and send only raw websocket native messages.
	// When only allow native messages is a fact?
	// When the registered namespace is just one and it's empty
	// and contains only one registered event which is the `OnNativeMessage`.
	// When `Events{...}` is used instead of `Namespaces{ "namespaceName": Events{...}}`
	// then the namespace is empty "".
	ws := websocket.New(
		websocket.DefaultGorillaUpgrader, websocket.Events{
			websocket.OnNativeMessage: func(nsConn *websocket.NSConn,
				msg websocket.Message) error {
				log.Printf("Server got: %s from [%s]", msg.Body, nsConn.Conn.ID())

				nsConn.Conn.Server().Broadcast(nsConn, msg)
				return nil
			},
		})
	ws.OnConnect = func(c *websocket.Conn) error {
		log.Printf("[%s] Connected to server!", c.ID())
		return nil
	}

	ws.OnDisconnect = func(c *websocket.Conn) {
		log.Printf("[%s] Disconnected from server", c.ID())
	}
	// register the server on an endpoint.
	// see the inline javascript code in the websockets.html,
	// this endpoint is used to connect to the server.
	app.Get("/echo", websocket.Handler(ws))
	// serve the javascript built'n client-side library,
	// see websockets.html script tags, this path is used.
	// app.Any("/iris-ws.js", websocket.ClientHandler())

}
