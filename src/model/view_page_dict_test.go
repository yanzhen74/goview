package model

import (
	"fmt"
	"testing"
)

func Test_get_view_page_dict(t *testing.T) {
	z, err := Read_view_page("../../config/resource/menu/WYG/RTM/SGYC/PK-CEH2.xml")
	if err != nil {
		fmt.Printf("error %v", err)
	}
	d := Get_view_page_dict(*z)
	for _, v := range *d {

		fmt.Println(*v)

	}
}
