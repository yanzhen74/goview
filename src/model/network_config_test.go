package model

import (
	"fmt"
	"testing"
)

func Test_read_network_config(t *testing.T) {
	z, err := Read_network_config("../../config/conf/NetWork.xml")
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	fmt.Println("Hello ", z)
	fmt.Println(z.NetWorkList[0].NetWorkProtocal)
}
