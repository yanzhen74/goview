package main

import (
	"net/http"
	_ "net/http/pprof"
	"github.com/yanzhen74/goview/src/routes"

	"github.com/kataras/iris"
	"github.com/yanzhen74/goview/src/controller"
	"github.com/yanzhen74/goview/src/goviewdb"
	"github.com/yanzhen74/goview/src/inits/parse"
	"github.com/yanzhen74/goview/src/model"
	// "github.com/yanzhen74/goview/src/routes"
)

func main() {
	app := iris.New()

	goviewdb.GwgDb = goviewdb.NewGWGDb("./db/gwg.db")

	parse.AppConfigParse()

	parse.DBSettingParse()
	goviewdb.MasterEngine().Sync2(new(model.User))
	goviewdb.SlaveEngine().Sync2(new(model.User))

	routes.Hub(app)

	app.HandleDir("/public", "./public")
	app.HandleDir("/config", "./config")
	app.HandleDir("/data", "./data")
	// 模板引擎采用html/template
	tmpl := iris.HTML("./views", ".html")

	// 在每个请求上 重新加载模板（开发模式）
	tmpl.Reload(true)
	app.RegisterView(tmpl)

	// init net
	controller.Init_network("config/conf/NetWork.xml")

	// Start processer for each frame in list
	controller.Init_0c_Processer("config/conf/ParameterDictionary.xml")

	// Receive network data
	controller.Run_network()

	go http.ListenAndServe(":10108", nil)
	// app.Run(iris.Addr(":"+parse.O.Other.Port), iris.WithoutServerError(iris.ErrServerClosed))
	app.Run(iris.Addr(":"+parse.AppConfig.Port), iris.WithConfiguration(iris.YAML("config/iris.yaml")))
}
