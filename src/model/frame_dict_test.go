package model

import (
	"fmt"
	"testing"
)

func Test_get_frame_dict_list(t *testing.T) {
	z, _ := Read_para_dict("../../config/conf/ParameterDictionary.xml")

	dicts := Get_frame_dict_list(z)
	for _, p := range *dicts {
		fmt.Println("hello, ", p)
	}
}
