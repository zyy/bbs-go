package weixin

import "github.com/chanxuehong/wechat/mp/core"

type MyHandler struct{}

func (MyHandler) ServeMsg(ctx *core.Context) {
	panic("implement me")
}
