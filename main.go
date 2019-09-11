package main

import (
	"github.com/kataras/iris"

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
	app.RegisterView(iris.HTML("./views", ".html"))

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

	app.Get("/", func(ctx iris.Context) {
		ctx.ViewData("menu", Menu{UserName: "oliver",
			Items: []ShipPackage{ShipPackage{Name: "YHYH", Children: []ShipPackage{ShipPackage{Name: "RTM", Children: nil}}}}})
		ctx.View("/index.html")
	})

	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}
