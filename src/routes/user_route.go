package routes

import (
	"strconv"
	"strings"
	"time"

	"github.com/yanzhen74/goview/src/model"
	"github.com/yanzhen74/goview/src/supports"
	"github.com/yanzhen74/goview/src/utils"

	"github.com/yanzhen74/goview/src/supports/vo"

	"github.com/yanzhen74/goview/src/middleware/jwts"

	"github.com/kataras/golog"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/hero"
)

func UserHub(party iris.Party) {
	var user = party.Party("/user")
	user.Post("/registe", hero.Handler(Registe))
	user.Post("/login", hero.Handler(Login))
	user.Post("/logout", hero.Handler(LoginOut))
	user.Get("/logout", func(ctx context.Context) { supports.Ok_(ctx, supports.RegisteSuccess) })
}

func Registe(ctx iris.Context) {
	var (
		err    error
		user   = new(model.User)
		effect int64
	)

	if err = ctx.ReadJSON(&user); err != nil {
		ctx.Application().Logger().Errorf("用户[%s]注册失败。%s", user.Username, err.Error())
		goto FAIL
	}

	user.Password = utils.AESEncrypt([]byte(user.Password))
	user.Enable = true
	user.CreateTime = time.Now()
	effect, err = model.CreateUser(user)
	if effect <= 0 || err != nil {
		ctx.Application().Logger().Errorf("用户[%s]注册失败。%s", user.Username, err.Error())
		goto FAIL
	}

	supports.Ok_(ctx, supports.RegisteSuccess)
	return

FAIL:
	supports.Error(ctx, iris.StatusInternalServerError, supports.RegisteFailur, nil)
	return
}

func LoginOut(ctx iris.Context) {
	supports.Ok(ctx, supports.LoginSuccess, nil)
	return
}

func Login(ctx iris.Context) {
	var (
		err        error
		user       = new(model.User)
		mUser      = new(model.User) // 查询数据后填充了数据的user
		ckPassword bool
		token      string
	)

	if err = ctx.ReadJSON(&user); err != nil {
		ctx.Application().Logger().Errorf("用户[%s]登录失败。%s", "", err.Error())
		supports.Error(ctx, iris.StatusBadRequest, supports.LoginFailur, nil)
		return
	}

	mUser.Username = user.Username
	has, err := model.GetUserByUsername(mUser)
	if err != nil {
		ctx.Application().Logger().Errorf("用户[%s]登录失败。%s", user.Username, err.Error())
		supports.Error(ctx, iris.StatusInternalServerError, supports.LoginFailur, nil)
		return
	}

	if !has { // 用户名不正确
		supports.Unauthorized(ctx, supports.UsernameFailur, nil)
		return
	}

	ckPassword = utils.CheckPWD(user.Password, mUser.Password)
	if !ckPassword {
		supports.Unauthorized(ctx, supports.PasswordFailur, nil)
		return
	}

	token, err = jwts.GenerateToken(mUser)
	golog.Infof("用户[%s], 登录生成token [%s]", mUser.Username, token)
	if err != nil {
		ctx.Application().Logger().Errorf("用户[%s]登录，生成token出错。%s", user.Username, err.Error())
		supports.Error(ctx, iris.StatusInternalServerError, supports.TokenCreateFailur, nil)
		return
	}

	supports.Ok(ctx, supports.LoginSuccess, vo.BuildUserVO(token, mUser))
	return
}

func UserDepTree(ctx iris.Context) {
	var (
		err  error
		deps []*model.Dep
	)

	if deps, err = model.GetAllDep(); err != nil {
		supports.Error(ctx, iris.StatusInternalServerError, supports.OptionFailur, nil)
		return
	}

	supports.Ok(ctx, supports.OptionSuccess, vo.BuildDepTreeForUser(deps))
	return
}

// 用户报表
func UserTable(ctx iris.Context) {
	var (
		err      error
		page     *supports.Pagination
		total    int64
		userList []*model.User
		//depId    int
	)

	//depId, err = ctx.URLParamInt("depId")
	page, err = supports.NewPagination(ctx)
	if err != nil {
		goto FAIL
	}

	userList, total, err = model.GetPaginationUsers(page)
	if err != nil {
		goto ERR
	}

	ctx.JSON(vo.TableVO{
		Total: total,
		Rows:  vo.BuildUserVOList(userList...),
	})
	return

FAIL:
	supports.Error(ctx, iris.StatusBadRequest, supports.ParseParamsFailur, nil)
	return
ERR:
	supports.Error(ctx, iris.StatusInternalServerError, supports.OptionFailur, nil)
	return
}

func UpdateUser(ctx iris.Context) {
	user := new(model.User)
	if err := ctx.ReadJSON(&user); err != nil {
		ctx.Application().Logger().Errorf("更新用户[%s]失败。%s", "", err.Error())
		supports.Error(ctx, iris.StatusBadRequest, supports.OptionFailur, nil)
		return
	}
	effect, err := model.UpdateUserById(user)
	if err != nil {
		ctx.Application().Logger().Errorf("更新用户[%s]失败。%s", "", err.Error())
		supports.Error(ctx, iris.StatusInternalServerError, supports.OptionFailur, nil)
		return
	}
	supports.Ok(ctx, supports.OptionSuccess, effect)
}

// 删除用户
func DeleteUsers(ctx iris.Context, uids string) {
	uidList := strings.Split(uids, ",")
	if len(uidList) == 0 {
		ctx.Application().Logger().Error("删除用户错误, 参数不对.")
		supports.Error(ctx, iris.StatusBadRequest, supports.ParseParamsFailur, nil)
		return
	}

	dUids := make([]int64, 0)
	for _, v := range uidList {
		if v == "" {
			continue
		}
		uid, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			ctx.Application().Logger().Errorf("删除用户错误, %s", err.Error())
			supports.Error(ctx, iris.StatusInternalServerError, supports.ParseParamsFailur, nil)
			return
		}
		dUids = append(dUids, uid)
	}

	effect, err := model.DeleteByUsers(dUids)
	if err != nil {
		ctx.Application().Logger().Errorf("删除用户错误, %s", err.Error())
		supports.Error(ctx, iris.StatusInternalServerError, supports.DeleteUsersFailur, nil)
		return
	}
	supports.Ok(ctx, supports.DeleteUsersSuccess, effect)
}
