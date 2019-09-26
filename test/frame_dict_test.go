package test

import (
	"fmt"
	"testing"

	"github.com/yanzhen74/goview/src/model"
)

func Test_get_frame_dict_list(t *testing.T) {
	z, _ := model.Read_para_dict("../config/conf/ParameterDictionary.xml")

	dicts := model.Get_frame_dict_list(z)
	for _, p := range *dicts {
		fmt.Println("hello, ", p)
	}
}
