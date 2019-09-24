package test

import (
	"fmt"
	"testing"

	"github.com/yanzhen74/goview/src/model"
)

func Test_read_network_config(t *testing.T) {
	z, err := model.Read_network_config("../config/conf/NetWork.xml")
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	fmt.Println("Hello ", z)
	fmt.Println(z.NetWorkList[0].NetWorkProtocal)
}
