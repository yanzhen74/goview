package controller

import (
	"log"
	"time"

	"github.com/kataras/iris/websocket"
)

func PublishPkg(nsConn *websocket.NSConn,
	msg websocket.Message, view_chan chan string) error {

	for i := 0; ; /*i < 100*/ i++ {
		// to be deleted, test if channel was blocked
		if i%100 == 0 {
			time.Sleep(time.Duration(2e9))
		}
		pkg, ok := <-view_chan
		if ok == false {
			log.Printf("channel has closed so publishPkg exit too\n")
			break
		}
		// log.Printf("receive %d\n", i)
		nsConn.Emit(string(msg.Body), []byte(pkg))
	}
	return nil
}
