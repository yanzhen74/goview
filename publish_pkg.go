package main

import (
	"fmt"
	"time"

	"github.com/kataras/iris/websocket"
)

func publishPkg(nsConn *websocket.NSConn,
	msg websocket.Message) error {

	for i := 0; i < 100; i++ {

		pkg := fmt.Sprintf("grid data %d", i)
		nsConn.Emit(string(msg.Body), []byte(pkg))
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}
