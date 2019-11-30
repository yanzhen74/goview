package controller

import (
	"fmt"
	"testing"

	"github.com/yanzhen74/goview/src/controller"
	"github.com/yanzhen74/goview/src/model"
)

// 每当return，k值不更新
//............................................
//process_type = ""
// expect output:
// 0xff 255 0xff true
// 0xff 255 0xff false
// expect output:
// 0xff 255 0xff true
// 0xff 255 0xff false
//............................................

func Test_param_Transfer_None(t *testing.T) {
	var para model.Para
	temp_strcode := "ff"
	temp_strresultvalue := "255"
	strCode := &temp_strcode
	strResultValue := &temp_strresultvalue
	temp_strresult := temp_strresultvalue
	strResult := &temp_strresult
	i := 0
	for i < 2 {
		if i == 0 {
			para.Process_type = ""
			para.Process_unit = "byte"
			para.Process_start = "0"
			para.Process_end = "1"
			Normal, Error := controller.Param_Transfer(para, strCode, strResult, strResultValue)
			fmt.Printf("%s %s %s %t %s\n", *strCode, *strResult, *strResultValue, Normal, Error)
		} else if i == 1 {
			para.Process_type = "aaa"
			para.Process_unit = "byte"
			para.Process_start = "0"
			para.Process_end = "1"
			Normal, Error := controller.Param_Transfer(para, strCode, strResult, strResultValue)
			fmt.Printf("%s %s %s %t %s\n", *strCode, *strResult, *strResultValue, Normal, Error)
		}
		i++
	}
}

//............................................
//process_type = "raw"; process_unit = "byte"
// expect output:
// 0x01 TestByte 0x01 true
// 0x02 TestByte 0x02 true
// 0x03 TestByte 0x03 true
// 0x04 TestByte 0x04 true
// 0x05 0x05 0x05 false
// expect output:
// 0x01 TestByte 0x01 true
// 0x02 TestByte 0x02 true
// 0x03 TestByte 0x03 true
// 0x04 TestByte 0x04 true
// 0x05 0x05 0x05 false
//............................................

func Test_param_Transfer_Byte(t *testing.T) {
	var para model.Para
	temp_strcode := ""
	temp_strresultvalue := ""
	strCode := &temp_strcode
	strResultValue := &temp_strresultvalue
	temp_strresult := temp_strresultvalue
	strResult := &temp_strresult
	para.Process_type = "raw"
	para.Process_unit = "byte"
	para.Process_start = "0"
	para.Process_end = "1"
	para.ParaRangeList = make([]model.ParaRange, 5, 5)
	i := 0
	for i < 5 {
		if i == 0 {
			temp_strcode = "01"
			para.ParaRangeList[i].ParaRangeSpecification = "TestByte"
			para.ParaRangeList[i].Alarm_min = "01"
			para.ParaRangeList[i].Alarm_min_equal = "true"
			para.ParaRangeList[i].Alarm_max = "01"
			para.ParaRangeList[i].Alarm_min_equal = "true"
			Normal, Error := controller.Param_Transfer(para, strCode, strResult, strResultValue)
			fmt.Printf("%s %s %s %t %s\n", *strCode, *strResult, *strResultValue, Normal, Error)
		} else if i == 1 {
			temp_strcode = "02"
			para.ParaRangeList[i].ParaRangeSpecification = "TestByte"
			para.ParaRangeList[i].Alarm_min = "02"
			para.ParaRangeList[i].Alarm_min_equal = "true"
			para.ParaRangeList[i].Alarm_max = "02"
			para.ParaRangeList[i].Alarm_min_equal = "true"
			Normal, Error := controller.Param_Transfer(para, strCode, strResult, strResultValue)
			fmt.Printf("%s %s %s %t %s\n", *strCode, *strResult, *strResultValue, Normal, Error)
		} else if i == 2 {
			temp_strcode = "03"
			para.ParaRangeList[i].ParaRangeSpecification = "TestByte"
			para.ParaRangeList[i].Alarm_min = "03"
			para.ParaRangeList[i].Alarm_min_equal = "true"
			para.ParaRangeList[i].Alarm_max = "03"
			para.ParaRangeList[i].Alarm_min_equal = "true"
			Normal, Error := controller.Param_Transfer(para, strCode, strResult, strResultValue)
			fmt.Printf("%s %s %s %t %s\n", *strCode, *strResult, *strResultValue, Normal, Error)
		} else if i == 3 {
			temp_strcode = "04"
			para.ParaRangeList[i].ParaRangeSpecification = "TestByte"
			para.ParaRangeList[i].Alarm_min = "04"
			para.ParaRangeList[i].Alarm_min_equal = "true"
			para.ParaRangeList[i].Alarm_max = "04"
			para.ParaRangeList[i].Alarm_min_equal = "true"
			Normal, Error := controller.Param_Transfer(para, strCode, strResult, strResultValue)
			fmt.Printf("%s %s %s %t %s\n", *strCode, *strResult, *strResultValue, Normal, Error)
		} else if i == 4 {
			temp_strcode = "05"
			para.ParaRangeList[i].ParaRangeSpecification = "TestByte"
			para.ParaRangeList[i].Alarm_min = "01"
			para.ParaRangeList[i].Alarm_min_equal = "true"
			para.ParaRangeList[i].Alarm_max = "01"
			para.ParaRangeList[i].Alarm_min_equal = "true"
			Normal, Error := controller.Param_Transfer(para, strCode, strResult, strResultValue)
			fmt.Printf("%s %s %s %t %s\n", *strCode, *strResult, *strResultValue, Normal, Error)
		}
		i++
	}
}

// ............................................
// process_type = "raw"; process_unit = "bit"
// expect output:
// 0xe TestBit 1... true
// 0xe 1110 1110 false
// 0xee TestBit 1....... true
// 0xee 11101110 11101110 false
// expect output:
// 0xe TestBit 1... true
// 0xe 1110 1110 false
// 0xee TestBit 1....... true
// 0xee 11101110 11101110 false
// ............................................

func Test_param_Transfer_Bit(t *testing.T) {
	var para model.Para
	temp_strcode := ""
	temp_strresultvalue := ""
	strCode := &temp_strcode
	strResultValue := &temp_strresultvalue
	temp_strresult := temp_strresultvalue
	strResult := &temp_strresult
	para.ParaRangeList = make([]model.ParaRange, 4, 4)
	i := 0
	for i < 4 {
		if i == 0 {
			temp_strcode = "1111"
			temp_strresultvalue = "15"
			para.Process_type = "raw"
			para.Process_unit = "bit"
			para.Process_start = "0"
			para.Process_end = "0"
			para.ParaRangeList[i].ParaRangeSpecification = "TestBit"
			para.ParaRangeList[i].Alarm_min = "0"
			para.ParaRangeList[i].Alarm_min_equal = "true"
			para.ParaRangeList[i].Alarm_max = "1"
			para.ParaRangeList[i].Alarm_max_equal = "true"
			Normal, Error := controller.Param_Transfer(para, strCode, strResult, strResultValue)
			fmt.Printf("%s %s %s %t %s\n", *strCode, *strResult, *strResultValue, Normal, Error)
		} else if i == 1 {
			temp_strcode = "1111"
			temp_strresultvalue = "15"
			para.Process_type = "raw"
			para.Process_unit = "bit"
			para.Process_start = "0"
			para.Process_end = "3"
			para.ParaRangeList[i].ParaRangeSpecification = "TestBit"
			para.ParaRangeList[i].Alarm_min = "0"
			para.ParaRangeList[i].Alarm_min_equal = "true"
			para.ParaRangeList[i].Alarm_max = "1"
			para.ParaRangeList[i].Alarm_max_equal = "true"
			Normal, Error := controller.Param_Transfer(para, strCode, strResult, strResultValue)
			fmt.Printf("%s %s %s %t %s\n", *strCode, *strResult, *strResultValue, Normal, Error)
		} else if i == 2 {
			temp_strcode = "11111111"
			temp_strresultvalue = "238"
			para.Process_type = "raw"
			para.Process_unit = "bit"
			para.Process_start = "0"
			para.Process_end = "0"
			para.ParaRangeList[i].ParaRangeSpecification = "TestBit"
			para.ParaRangeList[i].Alarm_min = "0"
			para.ParaRangeList[i].Alarm_min_equal = "true"
			para.ParaRangeList[i].Alarm_max = "1"
			para.ParaRangeList[i].Alarm_max_equal = "true"
			Normal, Error := controller.Param_Transfer(para, strCode, strResult, strResultValue)
			fmt.Printf("%s %s %s %t %s\n", *strCode, *strResult, *strResultValue, Normal, Error)
		} else if i == 3 {
			temp_strcode = "11111111"
			temp_strresultvalue = "238"
			para.Process_type = "raw"
			para.Process_unit = "bit"
			para.Process_start = "0"
			para.Process_end = "7"
			para.ParaRangeList[i].ParaRangeSpecification = "TestBit"
			para.ParaRangeList[i].Alarm_min = "0"
			para.ParaRangeList[i].Alarm_min_equal = "true"
			para.ParaRangeList[i].Alarm_max = "1"
			para.ParaRangeList[i].Alarm_max_equal = "true"
			Normal, Error := controller.Param_Transfer(para, strCode, strResult, strResultValue)
			fmt.Printf("%s %s %s %t %s\n", *strCode, *strResult, *strResultValue, Normal, Error)
		}
		i++
	}
}

// ............................................
// process_type = "raw"; process_unit = "code"
// expect output:
// 0xffee TestCode 0xffee true
// 0xffee 0xee 0xee true
// ffff 0xffff 0xff false
// expect output:
// 0xffee TestCode 0xffee true
// 0xffee 0xee 0xee true
// ffff 0xffff 0xff false
// ............................................

func Test_param_Transfer_Code(t *testing.T) {
	var para model.Para
	temp_strcode := ""
	temp_strresultvalue := ""
	strCode := &temp_strcode
	strResultValue := &temp_strresultvalue
	temp_strresult := temp_strresultvalue
	strResult := &temp_strresult
	para.ParaRangeList = make([]model.ParaRange, 3, 3)
	i := 0
	for i < 3 {
		if i == 0 {
			temp_strcode = "ffee"
			temp_strresultvalue = "ffee"
			para.Process_type = "raw"
			para.Process_unit = "code"
			para.Process_start = "0"
			para.Process_end = "3"
			para.ParaRangeList[i].ParaRangeSpecification = "TestCode"
			para.ParaRangeList[i].Alarm_min = "ffee"
			para.ParaRangeList[i].Alarm_min_equal = "true"
			para.ParaRangeList[i].Alarm_max = "ffee"
			para.ParaRangeList[i].Alarm_max_equal = "true"
			Normal, Error := controller.Param_Transfer(para, strCode, strResult, strResultValue)
			fmt.Printf("%s %s %s %t %s\n", *strCode, *strResult, *strResultValue, Normal, Error)
		} else if i == 1 {
			temp_strcode = "ffee"
			temp_strresultvalue = "ffee"
			para.Process_type = "raw"
			para.Process_unit = "code"
			para.Process_start = "2"
			para.Process_end = "3"
			para.ParaRangeList[i].ParaRangeSpecification = ""
			para.ParaRangeList[i].Alarm_min = "ee"
			para.ParaRangeList[i].Alarm_min_equal = "true"
			para.ParaRangeList[i].Alarm_max = "ee"
			para.ParaRangeList[i].Alarm_max_equal = "true"
			Normal, Error := controller.Param_Transfer(para, strCode, strResult, strResultValue)
			fmt.Printf("%s %s %s %t %s\n", *strCode, *strResult, *strResultValue, Normal, Error)
		} else if i == 2 {
			temp_strcode = "ffff"
			temp_strresultvalue = "ffff"
			para.Process_type = "raw"
			para.Process_unit = "code"
			para.Process_start = "0"
			para.Process_end = "1"
			para.ParaRangeList[i].ParaRangeSpecification = "TestCode"
			para.ParaRangeList[i].Alarm_min = "qqq"
			para.ParaRangeList[i].Alarm_min_equal = "true"
			para.ParaRangeList[i].Alarm_max = "qqq"
			para.ParaRangeList[i].Alarm_max_equal = "true"
			Normal, Error := controller.Param_Transfer(para, strCode, strResult, strResultValue)
			fmt.Printf("%s %s %s %t %s\n", *strCode, *strResult, *strResultValue, Normal, Error)
		}
		i++
	}
}

// ............................................
// process_type = "raw"; process_unit = "longcode"
// expect output:
// 0xaaaa TestLongCode 0xaaaa true
// 0xaaaa 0xaaaa 0xaaaa true
// expect output:
// 0xaaaa TestLongCode 0xaaaa true
// 0xaaaa 0xaaaa 0xaaaa true
// ............................................

func Test_param_Transfer_Longcode(t *testing.T) {
	var para model.Para
	temp_strcode := ""
	temp_strresultvalue := ""
	temp_strresult := temp_strresultvalue
	strCode := &temp_strcode
	strResultValue := &temp_strresultvalue
	strResult := &temp_strresult
	para.ParaRangeList = make([]model.ParaRange, 2, 2)

	for i := 0; i < 2; i++ {
		if i == 0 {
			temp_strcode = "aaaa"
			temp_strresultvalue = "aaaa"
			para.Process_type = "raw"
			para.Process_unit = "longcode"
			para.Process_start = "0"
			para.Process_end = "3"
			para.ParaRangeList[i].ParaRangeSpecification = "TestLongCode"
			para.ParaRangeList[i].Alarm_min = "ffff"
			para.ParaRangeList[i].Alarm_min_equal = "true"
			para.ParaRangeList[i].Alarm_max = "ffff"
			para.ParaRangeList[i].Alarm_max_equal = "true"
			Normal, Error := controller.Param_Transfer(para, strCode, strResult, strResultValue)
			fmt.Printf("%s %s %s %t %s\n", *strCode, *strResult, *strResultValue, Normal, Error)
		}
		if i == 1 {
			temp_strcode = "aaaa"
			temp_strresultvalue = "aaaa"
			para.Process_type = "raw"
			para.Process_unit = "longcode"
			para.Process_start = "0"
			para.Process_end = "3"
			para.ParaRangeList[i].ParaRangeSpecification = ""
			para.ParaRangeList[i].Alarm_min = "ffff"
			para.ParaRangeList[i].Alarm_min_equal = "true"
			para.ParaRangeList[i].Alarm_max = "ffff"
			para.ParaRangeList[i].Alarm_max_equal = "true"
			Normal, Error := controller.Param_Transfer(para, strCode, strResult, strResultValue)
			fmt.Printf("%s %s %s %t %s\n", *strCode, *strResult, *strResultValue, Normal, Error)
		}
	}
}

// ............................................
// process_type = "result"; Type = "int"
// expect output:
// 0xee TestInt 238 true
// 0xee 238 238 false
// expect output:
// 0xee TestInt 238 true
// 0xee 238 238 false
// ............................................

func Test_param_Transfer_Int(t *testing.T) {
	var para model.Para
	temp_strcode := ""
	temp_strresultvalue := ""
	temp_strresult := temp_strresultvalue
	strCode := &temp_strcode
	strResultValue := &temp_strresultvalue
	strResult := &temp_strresult
	para.ParaRangeList = make([]model.ParaRange, 2, 2)

	for i := 0; i < 2; i++ {
		if i == 0 {
			temp_strcode = "ee"
			temp_strresultvalue = "238"
			temp_strresult = "238"
			para.Process_type = "result"
			para.Type = "int"
			para.Process_start = "0"
			para.Process_end = "3"
			para.ParaRangeList[i].ParaRangeSpecification = "TestInt"
			para.ParaRangeList[i].Alarm_min = "10"
			para.ParaRangeList[i].Alarm_min_equal = "true"
			para.ParaRangeList[i].Alarm_max = "300"
			para.ParaRangeList[i].Alarm_max_equal = "true"
			Normal, Error := controller.Param_Transfer(para, strCode, strResult, strResultValue)
			fmt.Printf("%s %s %s %t %s\n", *strCode, *strResult, *strResultValue, Normal, Error)
		}
		if i == 1 {
			temp_strcode = "ee"
			temp_strresultvalue = "238"
			temp_strresult = "238"
			para.Process_type = "result"
			para.Type = "int"
			para.Process_start = "0"
			para.Process_end = "3"
			para.ParaRangeList[i].ParaRangeSpecification = "TestInt"
			para.ParaRangeList[i].Alarm_min = "10"
			para.ParaRangeList[i].Alarm_min_equal = "true"
			para.ParaRangeList[i].Alarm_max = "10"
			para.ParaRangeList[i].Alarm_max_equal = "true"
			Normal, Error := controller.Param_Transfer(para, strCode, strResult, strResultValue)
			fmt.Printf("%s %s %s %t %s\n", *strCode, *strResult, *strResultValue, Normal, Error)
		}
	}
}

// ............................................
// process_type = "result"; Type = "float"
// expect output:
// 0xee TestFloat 2.8 true
// 0xee 7.8 7.8 false
// expect output:
// 0xee TestFloat 2.8 true
// 0xee 7.8 7.8 false
// ............................................

func Test_param_Transfer_Float(t *testing.T) {
	var para model.Para
	temp_strcode := ""
	temp_strresultvalue := ""
	temp_strresult := temp_strresultvalue
	strCode := &temp_strcode
	strResultValue := &temp_strresultvalue
	strResult := &temp_strresult
	para.ParaRangeList = make([]model.ParaRange, 1, 2)

	for i := 0; i < 1; i++ {
		if i == 0 {
			temp_strcode = "ee"
			temp_strresultvalue = "2.8"
			temp_strresult = "2.8"
			para.Process_type = "result"
			para.Type = "float"
			para.Process_start = "0"
			para.Process_end = "3"
			para.ParaRangeList[i].ParaRangeSpecification = "TestFloat"
			para.ParaRangeList[i].Alarm_min = "0.0"
			para.ParaRangeList[i].Alarm_min_equal = "false"
			para.ParaRangeList[i].Alarm_max = "3.5"
			para.ParaRangeList[i].Alarm_max_equal = "false"
			Normal, Error := controller.Param_Transfer(para, strCode, strResult, strResultValue)
			fmt.Printf("%s %s %s %t %s\n", *strCode, *strResult, *strResultValue, Normal, Error)
		}
		if i == 0 {
			temp_strcode = "ee"
			temp_strresultvalue = "7.8"
			temp_strresult = "7.8"
			para.Process_type = "result"
			para.Type = "float"
			para.Process_start = "0"
			para.Process_end = "3"
			para.ParaRangeList[i].ParaRangeSpecification = "TestFloat"
			para.ParaRangeList[i].Alarm_min = "3.5"
			para.ParaRangeList[i].Alarm_min_equal = "false"
			para.ParaRangeList[i].Alarm_max = "6.8"
			para.ParaRangeList[i].Alarm_max_equal = "false"
			Normal, Error := controller.Param_Transfer(para, strCode, strResult, strResultValue)
			fmt.Printf("%s %s %s %t %s\n", *strCode, *strResult, *strResultValue, Normal, Error)
		}
	}
}
