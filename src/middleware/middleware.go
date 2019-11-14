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
	if checkURL(path) || checkURLMatch(path) {
		ctx.Next()
		return
	}

	if !jwts.Serve(ctx) {
		return
	}

	// Pass to real API
	ctx.Next()
}

func checkURLMatch(reqPath string) bool {
	resources := []string{"/gwgpic", "/echo", "/tab"}
	for _, path := range resources {
		if strings.Contains(reqPath, path) {
			return true
		}
	}
	return false
}

/**
return
	true:则跳过不需验证，如登录接口等...
	false:需要进一步验证
*/
func checkURL(reqPath string) bool {
	//config := iris.YAML("conf/app.yml")
	//ignoreURLs := config.GetOther()["ignoreURLs"].([]interface{})
	for _, v := range parse.AppConfig.IgnoreURLs {
		if reqPath == v {
			return true
		}
	}
	return false
}
