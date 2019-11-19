package api

import (
	"github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
)

type WeixinController struct {
	Ctx iris.Context
}

func (this *WeixinController) AnyCallback() {
	var (
		signature = this.Ctx.FormValue("signature")
		timestamp = this.Ctx.FormValue("timestamp")
		nonce     = this.Ctx.FormValue("nonce")
	)

	logrus.Info(signature, timestamp, nonce)

	if this.Ctx.Method() == "GET" { // GET请求，说明是微信来校验Token
		echostr := this.Ctx.FormValue("echostr")
		// boolean success = WeixinUtils.getWxMaService().checkSignature(timestamp, nonce, signature);
		// if (success) {
		// 	writeRspString(context, echostr);
		// } else {
		// 	writeRspString(context, "error");
		// }

		// TODO gaoyoubo @ 2019/11/19 这里要进行校验
		this.Ctx.WriteString(echostr)
	} else {
		body, err := this.Ctx.GetBody()
		logrus.Info(string(body), err)
	}
}
