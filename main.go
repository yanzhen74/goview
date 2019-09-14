package main

import (
	"fmt"

	"github.com/kataras/iris"
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
	tmpl := iris.HTML("./views", ".html")
	tmpl.Reload(true)
	app.RegisterView(tmpl)

	// template
	// t := template.New("ex")
	// t, _ = t.Parse(`hello {{.UserName}}!
	// 	{{range .Items}}
	// 		{{with .Children}}
	// 		{{range .}}
	// 			child name is {{.Name}}
	// 		{{end}}
	// 		{{end}}
	// 	{{end}}
	// 	`)
	// t, _ = t.ParseFiles("views/index.html")
	// m := menu{UserName: "oliver",
	// 	Items: []ShipPackage{ShipPackage{Name: "YHYH", Children: []ShipPackage{ShipPackage{Name: "RTM", Children: nil}}}}}
	// t.Execute(os.Stdout, m)
	pages, err := model.Init_pages("config/resource/menu")

	if err != nil {
		fmt.Printf("error is %v", err)
		return
	}

	app.Get("/", func(ctx iris.Context) {
		ctx.ViewData("menu", pages)
		ctx.View("/index.html")
	})

	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}
