package weixin

import (
	"github.com/chanxuehong/wechat/mp/core"
	"github.com/chanxuehong/wechat/mp/message/callback/request"

	"github.com/mlogclub/bbs-go/services"
)

type MyHandler struct{}

func (m MyHandler) ServeMsg(ctx *core.Context) {
	switch ctx.MixedMsg.EventType {
	case request.EventTypeSubscribe: // 关注事件, 包括点击关注和扫描二维码(公众号二维码和公众号带参数二维码)关注
		m.HandleSubscribe(ctx)
		break
	case request.EventTypeScan: // 已经关注的用户扫描带参数二维码事件
		m.HandleScan(ctx)
		break
	case request.EventTypeUnsubscribe: // 取消关注事件
		m.HandleOther(ctx)
		break
	case request.EventTypeLocation: // 上报地理位置事件
		m.HandleOther(ctx)
		break
	}
}

// 关注事件
func (MyHandler) HandleSubscribe(ctx *core.Context) {
	weixinLog.Info(string(ctx.MsgPlaintext))

	msg := request.GetSubscribeEvent(ctx.MixedMsg)
	var (
		openId     = msg.FromUserName
		sceneId, _ = msg.Scene()
	)
	if len(sceneId) > 0 {
		weixinLog.Info("自然关注量...", openId)
		return
	}

	loginQrcode := GetLoginQrcode(sceneId)
	if loginQrcode == nil {
		weixinLog.Error("微信登录，二维码超时")
		SetLoginError(sceneId, "登录失败，二维码超时")
	} else {
		thirdAccount, err := services.ThirdAccountService.GetOrCreateByWeixin(openId)
		if err != nil {
			weixinLog.Error("微信登录，处理第三方账号时错误", err)
			SetLoginError(sceneId, "登录失败：01")
		} else {
			user, codeErr := services.UserService.SignInByThirdAccount(thirdAccount)
			if codeErr != nil {
				weixinLog.Error("微信登录，处理用户信息时错误", codeErr)
				SetLoginError(sceneId, "登录失败：02")
			} else {
				weixinLog.Info("微信登录，成功：", openId)
				SetLoginSuccess(sceneId, user.Id)
			}
		}
	}
}

// 已经关注的用户扫描带参数二维码事件
func (MyHandler) HandleScan(ctx *core.Context) {
	// msg := request.GetScanEvent(ctx.MixedMsg)
	weixinLog.Info(string(ctx.MsgPlaintext))
}

// 其他
func (MyHandler) HandleOther(ctx *core.Context) {
	weixinLog.Info(string(ctx.MsgPlaintext)) // 打印消息
}
