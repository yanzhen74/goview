package controller

import (
	"fmt"
	"time"

	"github.com/kataras/iris/websocket"
)

func publishPkg(nsConn *websocket.NSConn,
	msg websocket.Message) error {

	for i := 0; i < 100; i++ {

		pkg := fmt.Sprintf("0,%d,%d,%d;1,%d,%d,%d;2,%d,%d,%d", i, i, i, i, i, i, i, i, i)
		nsConn.Emit(string(msg.Body), []byte(pkg))
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}
