package model

import (
	"fmt"
	"testing"
)

func Test_read_para_dict(t *testing.T) {
	z, err := Read_para_dict("../../config/conf/ParameterDictionary.xml")
	if err != nil {
		fmt.Printf("error %v", err)
	}
	fmt.Println("Hello ", z)
	fmt.Println(z.AircraftList[0].DataTypeList[0].PayLoadList[0].SubAddressList[0].ParaList[0].Name)
}
