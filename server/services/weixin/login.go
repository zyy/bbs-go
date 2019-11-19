package weixin

import (
	"time"

	"github.com/chanxuehong/wechat/mp/qrcode"
	"github.com/goburrow/cache"
	"github.com/mlogclub/simple"
)

type LoginQrcode struct {
	TempQrcode *qrcode.TempQrcode // 二维码
	Status     int                // 是否登录成功, 0: 进行中, 1: 成功, 2: 失败
	UserId     int64              // 登录成功后用户编号
	ErrorMsg   string             // 错误消息
}

// 登录二维码，key：sceneId，value：LoginQrcode
var qrcodeSceneIdCache = cache.NewLoadingCache(
	func(key cache.Key) (value cache.Value, e error) {
		tempQrcode, err := qrcode.CreateStrSceneTempQrcode(GetClient(), key.(string), 60*30)
		if err != nil {
			e = err
			return
		}
		value = &LoginQrcode{
			TempQrcode: tempQrcode,
		}
		return
	},
	cache.WithMaximumSize(1000),
	cache.WithExpireAfterAccess(30*time.Minute),
)

// 创建登录二维码
func GenerateLoginQrcode() (sceneId string, loginQrcode *LoginQrcode, err error) {
	sceneId = simple.Uuid()
	val, e := qrcodeSceneIdCache.Get(sceneId)
	if e != nil {
		err = e
		return
	}
	loginQrcode = val.(*LoginQrcode)
	return
}

// 创建登录二维码
func GetLoginQrcode(sceneId string) *LoginQrcode {
	val, found := qrcodeSceneIdCache.GetIfPresent(sceneId)
	if !found {
		return nil
	}
	return val.(*LoginQrcode)
}

// 登录失败
func SetLoginError(sceneId, errorMsg string) {
	loginQrcode := GetLoginQrcode(sceneId)
	if loginQrcode == nil {
		loginQrcode = &LoginQrcode{}
	}
	loginQrcode.Status = 2
	loginQrcode.ErrorMsg = errorMsg
	qrcodeSceneIdCache.Put(sceneId, loginQrcode)
}

// 登录成功
func SetLoginSuccess(sceneId string, userId int64) {
	loginQrcode := GetLoginQrcode(sceneId)
	if loginQrcode == nil {
		loginQrcode = &LoginQrcode{}
	}
	loginQrcode.Status = 1
	loginQrcode.UserId = userId
	qrcodeSceneIdCache.Put(sceneId, loginQrcode)
}
