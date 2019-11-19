package weixin

import (
	"sync"

	"github.com/chanxuehong/wechat/mp/core"
	"github.com/sirupsen/logrus"
)

const (
	AppId          = "wx269fff3c0c082b0d"
	AppSecret      = "c64ac7319513f7fa6168114e2f06ca3b"
	Token          = "IcoRXFWVkZANT7aKXEl6HFL9ECgZ4Ohz"
	EncodingAESKey = "7tCEWkYJHnYwAdcyAP2sy1Dv5pGdfb4SRvAPUvfWYWI"
)

var weixinLog = logrus.WithFields(logrus.Fields{
	"weixin": true,
	"appId":  AppId,
})

var (
	clientOnce sync.Once
	serverOnce sync.Once

	ats    core.AccessTokenServer
	client *core.Client
	serve  *core.Server
)

func GetClient() *core.Client {
	clientOnce.Do(func() {
		ats = core.NewDefaultAccessTokenServer(AppId, AppSecret, nil)
		client = core.NewClient(ats, nil)
	})
	return client
}

func GetServer() *core.Server {
	serverOnce.Do(func() {
		msgHandler := MyHandler{}
		errorHandler := MyErrorHandler{}
		serve = core.NewServer("", AppId, Token, EncodingAESKey, msgHandler, errorHandler)
	})
	return serve
}
