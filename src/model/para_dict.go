package model

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type Aircrafts struct {
	XMLName      xml.Name   `xml:"AirCrafts"`
	AircraftList []Aircraft `xml:"AirCraft"`
}

type Aircraft struct {
	XMLName      xml.Name   `xml:"AirCraft"`
	Name         string     `xml:"Name,attr"`
	DataTypeList []DataType `xml:"DataType"`
}

type DataType struct {
	XMLNode     xml.Name  `xml:"DataType"`
	Name        string    `xml:"Name,attr"`
	PayLoadList []PayLoad `xml:"PayLoad"`
}

type PayLoad struct {
	XMLNode        xml.Name     `xml:"PayLoad"`
	Name           string       `xml:"Name,attr"`
	SubAddressList []SubAddress `xml:"SubAddress"`
}

type SubAddress struct {
	XMLNode  xml.Name `xml:"SubAddress"`
	Name     string   `xml:"Name,attr"`
	ParaList []Para   `xml:"Para"`
}

type Para struct {
	XMLNode       xml.Name    `xml:"Para"`
	Name          string      `xml:"name,attr"`
	ID            string      `xml:"id,attr"`
	ParaKey       string      `xml:"Parakey,attr"`
	Type          string      `xml:"type,attr"`
	Unit          string      `xml:"unit,attr"`
	Process_type  string      `xml:"process_type,attr"`
	Process_unit  string      `xml:"process_unit,attr"`
	Process_start string      `xml:"process_start,attr"`
	Process_end   string      `xml:"process_end,attr"`
	ParaRangeList []ParaRange `xml:"ParaRange"`
}

type ParaRange struct {
	XMLNode                xml.Name `xml:"ParaRange"`
	RangeSeq               string   `xml:"RangeSeq,attr"`
	ParaRangeSpecification string   `xml:"ParaRangeSpecification,attr"`
	Alarm_min              string   `xml:"alarm_min,attr"`
	Alarm_min_equal        string   `xml:"alarm_min_equal,attr"`
	Alarm_max              string   `xml:"alarm_max,attr"`
	Alarm_max_equal        string   `xml:"alarm_max_equal,attr"`
}

func Read_para_dict(filename string) (*Aircrafts, error) {
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

	v := Aircrafts{}
	err = xml.Unmarshal(data, &v)
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil, err
	}

	fmt.Println(v)

	return &v, err
}
