package routes

import (
	"fmt"

	"github.com/kataras/iris"
	"github.com/yanzhen74/goview/src/model"
)

func HomeHub(party iris.Party) {
	home := party.Party("/")

	pages, err := model.Init_pages("config/resource/menu")
	if err != nil {
		fmt.Printf("error is %v", err)
		return
	}

	home.Get("/", func(ctx iris.Context) { // 首页模块
		//username, password, _ := ctx.Request().BasicAuth()
		//log.Printf("%s, %s\n", username, password)
		ctx.ViewData("menu", pages)
		ctx.View("index.html")
	})
}
