package model

type ViewType struct {
	MissionID      string
	DataType       string
	PayloadName    string
	SubAddressName string
}

type ViewDict struct {
	View_type ViewType
	ParaList  []Para_Page
}

func Get_view_page_dict(view Paras) *[]*ViewDict {
	var viewdicts = make([]*ViewDict, 0, 10)
	for _, a := range view.ParaList {
		bNewDict := true
		for _, d := range viewdicts {
			if a.SubAddressName == d.View_type.SubAddressName && a.PayloadName == d.View_type.PayloadName {
				d.ParaList = append(d.ParaList, a)
				bNewDict = false
				break
			}
		}

		if bNewDict {
			dict := new(ViewDict)
			dict.View_type = ViewType{
				view.MissionID, view.DataType, a.PayloadName, a.SubAddressName}
			dict.ParaList = make([]Para_Page, 0, 10)
			dict.ParaList = append(dict.ParaList, a)
			viewdicts = append(viewdicts, dict)
		}
	}
	return &viewdicts
}
