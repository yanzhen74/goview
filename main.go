package main

import (
	"fmt"

	"github.com/kataras/iris"
	"github.com/yanzhen74/goview/src/controller"
	"github.com/yanzhen74/goview/src/model"

	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

type ShipPackage struct {
	Name     string
	Children []ShipPackage
}
type Menu struct {
	UserName string
	Items    []ShipPackage
}

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")
	app.Use(recover.New())
	app.Use(logger.New())
	app.HandleDir("/public", "./public")
	app.HandleDir("/config", "./config")
	// 模板引擎采用html/template
	tmpl := iris.HTML("./views", ".html")
	// 在每个请求上 重新加载模板（开发模式）
	tmpl.Reload(true)
	app.RegisterView(tmpl)

	// map frame page
	controller.Frame_page_map = make(map[string]int)

	pages, err := model.Init_pages("config/resource/menu")

	if err != nil {
		fmt.Printf("error is %v", err)
		return
	}

	// map
	z, _ := model.Read_para_dict("config/conf/ParameterDictionary.xml")

	controller.Dicts = model.Get_frame_dict_list(z)

	controller.Frame_page_map[pages.Branches[0].Branches[0].Branches[0].Branches[0].Curdir] = 0
	controller.Frame_page_map[pages.Branches[1].Branches[0].Branches[0].Branches[0].Curdir] = 1
	go controller.Process0cPkg((*controller.Dicts)[0])
	go controller.Process0cPkg((*controller.Dicts)[1])

	// websocket
	controller.SetupWebsocket(app)

	app.Get("/tab/{page:path}", func(ctx iris.Context) {
		page := ctx.Params().Get("page")
		fmt.Println("open:", page)
		paras, err := model.Read_view_page(page)
		if err != nil {
			fmt.Printf("error is %v", err)
			return
		}
		fmt.Println("paras: ", paras)
		ctx.ViewData("paras", paras)
		ctx.View("/tab.html")
	})

	app.Get("/", func(ctx iris.Context) {
		ctx.ViewData("menu", pages)
		ctx.View("/index.html")
	})
	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}
