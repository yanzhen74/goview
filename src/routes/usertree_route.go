package routes

import (
	"fmt"
	"model"

	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/yanzhen74/goview/src/controller"
	"github.com/yanzhen74/goview/src/middleware/jwts"
	"github.com/yanzhen74/goview/src/supports"
)

func UserTreeHub(party iris.Party) {
	var usertree = party.Party("/usertree")
	usertree.Get("/add", hero.Handler(UserTreeAdd))
	usertree.Get("/edit", hero.Handler(UserTreeEdit))
	usertree.Get("/del", hero.Handler(UserTreeDel))
	usertree.Post("/keysearch", hero.Handler(UserTreeKeySearch))
}

func UserTreeAdd(ctx iris.Context) {
	claims := jwts.GetJWTClaims(ctx)
	if claims["username"] == "GUEST" {
		fmt.Println("ok")
	}

	// role := claims["role"]

	ctx.View("/usertree/add.html")
}

func UserTreeEdit(ctx iris.Context) {

}

func UserTreeDel(ctx iris.Context) {

}

func UserTreeKeySearch(ctx iris.Context) {
	type key_search struct {
		Key            string `json:"key_search"`
		MissionID      string `json:"mission_id"`
		DataType       string `json:"data_type"`
		PayloadName    string `json:"payload_name"`
		SubAddressName string `json:"subaddress_name"`
	}

	var (
		err error
		key = new(key_search)
	)

	if err = ctx.ReadJSON(&key); err != nil {
		ctx.Application().Logger().Errorf("search [%s] not valid. %s", key.Key, err.Error())
		supports.Error(ctx, iris.StatusBadRequest, "not valid", nil)
	}

	viewdict := new(model.ViewDict)
	viewdict.View_type.MissionID = key.MissionID
	viewdict.View_type.DataType = key.DataType
	viewdict.View_type.PayloadName = key.PayloadName
	viewdict.View_type.SubAddressName = key.SubAddressName
	viewdict.ParaList = make([]model.Para_Page, 0, 10)
	for _, d := range *controller.Dicts {
		if key.MissionID == d.Frame_type.MissionID &&
			key.DataType == d.Frame_type.DataType &&
			key.PayloadName == d.Frame_type.PayloadName &&
			key.SubAddressName == d.Frame_type.SubAddressName {
			for _, p := range d.ParaList {
				var para model.Para_Page
				para.XMLName = p.XMLNode
				para.Name = p.Name
				para.ID = string(p.ID)
				para.ParaKey = p.ParaKey
				para.Type = p.Type
				para.Unit = p.Unit
				para.PayloadName = key.PayloadName
				viewdict.ParaList = append(viewdict.ParaList, para)
			}
		}
	}
	supports.Ok(ctx, "ok", viewdict)

}
