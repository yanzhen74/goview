package routes

import (
	"fmt"
	"strconv"
	"github.com/yanzhen74/goview/src/supports"

	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/yanzhen74/goview/src/goviewdb"
)

func GwgPicHub(party iris.Party) {
	gwgpic := party.Party("/gwgpic")

	gwgpic.Get("/", hero.Handler(ListGwgPic))

	gwgpic.Get("/{pic_id}", hero.Handler(ViewGwgPic))

	gwgpic.Get("/list.dd", hero.Handler(FetchGwgPicList))

	gwgpic.Get("/save", hero.Handler(SaveGwgPic))
}

func SaveGwgPic(ctx iris.Context) {
	url := ctx.URLParam("url")

	goviewdb.GwgDb.SavePic(url)

	supports.Ok(ctx, supports.LoginSuccess, url)
}

func FetchGwgPicList(ctx iris.Context) {
	limit, _ := strconv.Atoi(ctx.URLParam("limit"))
	pageno, _ := strconv.Atoi(ctx.URLParam("page"))
	codeLike := ctx.URLParam("codeLike")
	// page := make([]goviewdb.GWGPic, 0, 20)
	page := goviewdb.GwgDb.ListPic(limit, pageno, codeLike)

	count := goviewdb.GwgDb.CountPic()

	// var i int64
	// for i = int64((pageno - 1) * limit); i < int64(pageno*limit); i++ {
	// 	page = append(page, goviewdb.GWGPic{Id: i, Url: "/data/gwg/1_17.jpeg", Camera: 0, Size: 0, ImageNo: 0, Time: "20H03M22S.333MS", CreatedAt: time.Now()})
	// }
	supports.Ok_page(ctx, supports.LoginSuccess, int(count), page)
}

func ListGwgPic(ctx iris.Context) {
	limit, _ := strconv.Atoi(ctx.URLParam("limit"))
	pageno, _ := strconv.Atoi(ctx.URLParam("page"))
	// page := make([]goviewdb.GWGPic, 0, 20)
	codeLike := ctx.URLParam("codeLike")
	piclist := goviewdb.GwgDb.ListPic(limit, pageno, codeLike)

	// show talbe of pics
	ctx.ViewData("piclist", piclist)
	ctx.View("gwgpic_list.html")

}

func ViewGwgPic(ctx iris.Context) {
	picNo := ctx.Params().Get("pic_id")
	fmt.Println("open:", picNo)

	// show table of paras
	ctx.ViewData("picNo", "/data/gwg/"+picNo)
	ctx.View("/gwgpic_view.html")
}
