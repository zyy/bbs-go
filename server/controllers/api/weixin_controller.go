package api

import (
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/services/weixin"
)

type WeixinController struct {
	Ctx iris.Context
}

func (this *WeixinController) AnyCallback() {
	// var (
	// 	signature = this.Ctx.FormValue("signature")
	// 	timestamp = this.Ctx.FormValue("timestamp")
	// 	nonce     = this.Ctx.FormValue("nonce")
	// )
	//
	// logrus.Info(signature, timestamp, nonce)

	// if this.Ctx.Method() == "GET" { // GET请求，说明是微信来校验Token
	// 	echostr := this.Ctx.FormValue("echostr")
	// 	// boolean success = WeixinUtils.getWxMaService().checkSignature(timestamp, nonce, signature);
	// 	// if (success) {
	// 	// 	writeRspString(context, echostr);
	// 	// } else {
	// 	// 	writeRspString(context, "error");
	// 	// }
	//
	// 	// TODO gaoyoubo @ 2019/11/19 这里要进行校验
	// 	_, _ = this.Ctx.WriteString(echostr)
	// } else {
	// 	body, err := this.Ctx.GetBody()
	// 	logrus.Info(string(body), err)
	// }

	weixin.GetServer().ServeHTTP(this.Ctx.ResponseWriter(), this.Ctx.Request(), nil)

}

// 登录二维码
func (this *WeixinController) GetLoginQrcode() *simple.JsonResult {
	sceneId, qr, err := weixin.GenerateLoginQrcode()
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.NewEmptyRspBuilder().Put("sceneId", sceneId).
		Put("url", qr.URL).
		JsonResult()
}
