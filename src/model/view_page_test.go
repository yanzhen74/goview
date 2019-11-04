package model

import (
	"fmt"
	"testing"
)

func Test_read_view_page(t *testing.T) {
	z, err := Read_view_page("../../config/resource/menu/WYG/RTM/SGYC/PK-CEH2.xml")
	if err != nil {
		fmt.Printf("error %v", err)
	}
	fmt.Println("hello ", z)
	// fmt.Println(z.Pages[0].)
}
