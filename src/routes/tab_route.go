package routes

import (
	"fmt"

	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/yanzhen74/goview/src/model"
)

// frontend view page paras model
var File_paras_map map[string]*model.Paras

func TabHub(party iris.Party) {
	tab := party.Party("/tab")

	// map file and paras - left menu
	File_paras_map = make(map[string]*model.Paras)

	tab.Get("/{page:path}", hero.Handler(OpenTab))
}

func OpenTab(ctx iris.Context) {
	page := ctx.Params().Get("page")
	fmt.Println("open:", page)

	// read view page paras from config file
	paras, err := model.Read_view_page(page)
	if err != nil {
		fmt.Printf("error is %v", err)
		return
	}

	// for websocket to use paras
	File_paras_map[paras.File] = paras

	// show table of paras
	ctx.ViewData("paras", paras)
	ctx.View("/tab.html")
}
