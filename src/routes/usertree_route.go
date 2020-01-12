package routes

import (
	"fmt"

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
		Key string `json:"key_search"`
	}

	var (
		err error
		key = new(key_search)
	)

	if err = ctx.ReadJSON(&key); err != nil {
		ctx.Application().Logger().Errorf("search [%s] not valid. %s", key.Key, err.Error())
		supports.Error(ctx, iris.StatusBadRequest, "not valid", nil)
	}

	supports.Ok(ctx, "ok", controller.Dicts)

}
