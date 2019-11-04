package middleware

import (
	"strings"

	"github.com/kataras/iris/context"
	"github.com/yanzhen74/goview/src/inits/parse"
	"github.com/yanzhen74/goview/src/middleware/jwts"
)

type Middleware struct {
}

func ServeHTTP(ctx context.Context) {
	path := ctx.Path()
	// 过滤静态资源、login接口、首页等...不需要验证
	if checkURL(path) || strings.Contains(path, "/static2") || strings.Contains(path, "echo") || strings.Contains(path, "tab") {
		ctx.Next()
		return
	}

	if !jwts.Serve(ctx) {
		return
	}

	// Pass to real API
	ctx.Next()
}

/**
return
	true:则跳过不需验证，如登录接口等...
	false:需要进一步验证
*/
func checkURL(reqPath string) bool {
	//config := iris.YAML("conf/app.yml")
	//ignoreURLs := config.GetOther()["ignoreURLs"].([]interface{})
	for _, v := range parse.O.Other.IgnoreURLs {
		if reqPath == v {
			return true
		}
	}
	return false
}
