package model

type FrameType struct {
	MissionID      string
	DataType       string
	PayloadName    string
	SubAddressName string
	ChanViewList   []chan string
	ChanViewReg    chan *View_page_regist_info
	ID             string
}

type FrameDict struct {
	Frame_type FrameType
	ParaList   []Para
}

func Get_frame_dict_list(aircraft Aircrafts) *[]FrameDict {
	var framedicts = make([]FrameDict, 0, 10)
	for _, a := range aircraft.AircraftList {
		for _, d := range a.DataTypeList {
			for _, p := range d.PayLoadList {
				for _, s := range p.SubAddressList {
					framedict := FrameDict{}
					framedict.Frame_type.MissionID = a.Name
					framedict.Frame_type.DataType = d.Name
					framedict.Frame_type.PayloadName = p.Name
					framedict.Frame_type.SubAddressName = s.Name
					framedict.Frame_type.ID = s.ID
					framedict.Frame_type.ChanViewList = make([]chan string, 0, 100)
					framedict.Frame_type.ChanViewReg = make(chan *View_page_regist_info, 10)
					framedict.ParaList = s.ParaList

					framedicts = append(framedicts, framedict)
				}
			}
		}
	}
	return &framedicts
}
