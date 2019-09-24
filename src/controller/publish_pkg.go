package controller

import (
	"github.com/kataras/iris/websocket"
)

func publishPkg(nsConn *websocket.NSConn,
	msg websocket.Message, view_chan chan string) error {

	for i := 0; i < 100; i++ {

		pkg := <-view_chan
		// pkg := fmt.Sprintf("0,%d,中国字%d,%d;1,%d,%d,%d;2,%d,%d,%d", i, i, i, i, i, i, i, i, i)
		nsConn.Emit(string(msg.Body), []byte(pkg))
		// time.Sleep(100 * time.Millisecond)
	}
	return nil
}
