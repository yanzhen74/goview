package vo

import "github.com/yanzhen74/goview/src/model"

// 前端需要的数据结构
type UserVO struct {
	*model.User
	Token string `json:"token"`
}

// 携带token
func BuildUserVO(token string, user *model.User) (uVO *UserVO) {
	uVO = &UserVO{
		user,
		token,
	}
	return
}

// 用户列表，不带token
func BuildUserVOList(userList ...*model.User) (userVOList []*UserVO) {
	for _, user := range userList {
		userVOList = append(userVOList, BuildUserVO("", user))
	}
	return
}
