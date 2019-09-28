package controller

import (
	"log"

	"github.com/kataras/iris/websocket"
)

func publishPkg(nsConn *websocket.NSConn,
	msg websocket.Message, view_chan chan string) error {

	for i := 0; ; /*i < 100*/ i++ {
		pkg, ok := <-view_chan
		if ok == false {
			log.Printf("channel has closed so publishPkg exit too\n")
			break
		}
		nsConn.Emit(string(msg.Body), []byte(pkg))
	}
	return nil
}
