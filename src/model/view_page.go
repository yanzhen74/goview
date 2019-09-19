package model

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type Paras struct {
	XMLName  xml.Name `xml:"Paras"`
	File     string
	ParaList []Para_Page `xml:"Para"`
}

type Para_Page struct {
	XMLName        xml.Name    `xml:"Para"`
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
	fmt.Println(string(data))
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

	fmt.Println(v)

	return &v, err
}
