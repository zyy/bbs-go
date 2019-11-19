package weixin

import (
	"sync"
	"time"

	"github.com/chanxuehong/wechat/mp/core"
	"github.com/chanxuehong/wechat/mp/qrcode"
	"github.com/goburrow/cache"
	"github.com/mlogclub/simple"
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

// 登录二维码，key：sceneId，value：qrcode.TempQrcode
var qrcodeSceneIdCache = cache.NewLoadingCache(
	func(key cache.Key) (value cache.Value, e error) {
		sceneId := key.(string)
		return qrcode.CreateStrSceneTempQrcode(GetClient(), sceneId, 60*30)
	},
	cache.WithMaximumSize(1000),
	cache.WithExpireAfterAccess(30*time.Minute),
)

// 创建登录二维码
func GenerateLoginQrcode() (sceneId string, qr *qrcode.TempQrcode, err error) {
	sceneId = simple.Uuid()
	val, e := qrcodeSceneIdCache.Get(sceneId)
	if e != nil {
		err = e
		return
	}
	qr = val.(*qrcode.TempQrcode)
	return
}

// 创建登录二维码
func GetLoginQrcode(sceneId string) (*qrcode.TempQrcode, error) {
	val, err := qrcodeSceneIdCache.Get(sceneId)
	if err != nil {
		return nil, err
	}
	return val.(*qrcode.TempQrcode), nil
}
