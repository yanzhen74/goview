package test

import (
	"fmt"
	"testing"

	"github.com/yanzhen74/goview/src/model"
)

func Test_read_view_page(t *testing.T) {
	z, err := model.Read_view_page("../config/resource/menu/WYG/RTM/SGYC/PK-CEH2.xml")
	if err != nil {
		fmt.Printf("error %v", err)
	}
	fmt.Println("hello ", z)
	// fmt.Println(z.Pages[0].)
}
