package main

import (
	"fmt"

	"github.com/kataras/iris"
	"github.com/yanzhen74/goview/src/controller"
	"github.com/yanzhen74/goview/src/inits/parse"
	"github.com/yanzhen74/goview/src/model"
	"github.com/yanzhen74/goview/src/routes"
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
	routes.Hub(app)
	parse.AppOtherParse()

	app.HandleDir("/public", "./public")
	app.HandleDir("/config", "./config")
	// 模板引擎采用html/template
	tmpl := iris.HTML("./views", ".html")

	// 在每个请求上 重新加载模板（开发模式）
	tmpl.Reload(true)
	app.RegisterView(tmpl)

	// init net config
	netConfig, err := model.Read_network_config("config/conf/NetWork.xml")
	if err != nil {
		fmt.Printf("error is %v", err)
		return
	}

	// init net
	controller.Init_network(netConfig)

	// map file and paras - left menu
	controller.File_paras_map = map[string]*model.Paras{}
	// Start receiver for each frame in list
	z, _ := model.Read_para_dict("config/conf/ParameterDictionary.xml")
	controller.Dicts = model.Get_frame_dict_list(z)

	for _, d := range *controller.Dicts {
		controller.Bind_network(d.Frame_type)
		go controller.Process0cPkg(d)
	}

	// Receive network data
	controller.Run_network()

	// websocket
	controller.SetupWebsocket(app)

	app.Run(iris.Addr(":"+parse.O.Other.Port), iris.WithoutServerError(iris.ErrServerClosed))
}
