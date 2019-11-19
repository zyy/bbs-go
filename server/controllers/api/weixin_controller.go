package api

import (
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/controllers/render"
	"github.com/mlogclub/bbs-go/services"
	"github.com/mlogclub/bbs-go/services/weixin"
)

type WeixinController struct {
	Ctx iris.Context
}

func (this *WeixinController) AnyCallback() {
	weixin.GetServer().ServeHTTP(this.Ctx.ResponseWriter(), this.Ctx.Request(), nil)
}

// 登录二维码
func (this *WeixinController) GetLoginQrcode() *simple.JsonResult {
	sceneId, loginQrcode, err := weixin.GenerateLoginQrcode()
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.NewEmptyRspBuilder().Put("sceneId", sceneId).
		Put("url", loginQrcode.TempQrcode.URL).
		JsonResult()
}

func (this *WeixinController) GetLoginStatus() *simple.JsonResult {
	sceneId := this.Ctx.FormValue("sceneId")
	if len(sceneId) == 0 {
		return simple.JsonErrorMsg("数据错误")
	}
	loginQrcode := weixin.GetLoginQrcode(sceneId)
	if loginQrcode == nil {
		return simple.JsonErrorMsg("数据超时")
	}
	if loginQrcode.Status == 0 { // 正在进行中
		return simple.NewEmptyRspBuilder().Put("status", loginQrcode.Status).JsonResult()
	} else if loginQrcode.Status == 1 { // 登录成功
		return this.GenerateLoginResult(loginQrcode.UserId, loginQrcode.Status)
	} else { // 登录失败
		return simple.JsonErrorMsg(loginQrcode.ErrorMsg)
	}
}

// user: login user
func (this *WeixinController) GenerateLoginResult(userId int64, status int) *simple.JsonResult {
	user := services.UserService.Get(userId)
	token, err := services.UserTokenService.Generate(userId)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.NewEmptyRspBuilder().
		Put("status", status).
		Put("token", token).
		Put("user", render.BuildUser(user)).
		JsonResult()
}
