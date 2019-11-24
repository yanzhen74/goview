package routes

import (
	"fmt"

	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/yanzhen74/goview/src/middleware/jwts"
)

func UserTreeHub(party iris.Party) {
	var usertree = party.Party("/usertree")
	usertree.Get("/add", hero.Handler(UserTreeAdd))
	usertree.Get("/edit", hero.Handler(UserTreeEdit))
	usertree.Get("/del", hero.Handler(UserTreeDel))
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
