package model

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type Paras struct {
	XMLName   xml.Name `xml:"Paras"`
	File      string
	MissionID string      `xml:"MissionID,attr"`
	DataType  string      `xml:"DataType,attr"`
	PageType  string      `xml:"PageType,attr"`
	ParaList  []Para_Page `xml:"Para"`
}

type Para_Page struct {
	XMLName        xml.Name `xml:"Para"`
	Index          string
	Name           string      `xml:"name,attr"`
	ID             string      `xml:"id,attr"`
	ParaKey        string      `xml:"Parakey,attr"`
	Type           string      `xml:"type,attr"`
	Unit           string      `xml:"unit,attr"`
	Process_type   string      `xml:"process_type,attr"`
	Process_unit   string      `xml:"process_unit,attr"`
	Process_start  string      `xml:"process_start,attr"`
	Process_end    string      `xml:"process_end,attr"`
	SubAddressName string      `xml:"SubAddressName,attr"`
	PayloadName    string      `xml:"PayloadName,attr"`
	Visible        string      `xml:"Visible,attr"`
	ParaRangeList  []ParaRange `xml:"ParaRange"`
}

func Read_view_page(filename string) (*Paras, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil, err
	}
	defer file.Close()
	data, err := ioutil.ReadFile(filename)
	// fmt.Println(string(data))
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil, err
	}

	v := Paras{}
	v.File = filename
	err = xml.Unmarshal(data, &v)
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil, err
	}

	// every param has only one Index but maybe multi ID
	index := 0
	for i, _ := range v.ParaList {
		v.ParaList[i].Index = strconv.Itoa(index)

		// fmt.Printf("index %d is %s\n", i, v.ParaList[i].Index)
		index += 1
	}
	//fmt.Println(v)

	return &v, err
}
