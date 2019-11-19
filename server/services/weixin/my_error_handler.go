package weixin

import (
	"net/http"
)

type MyErrorHandler struct{}

func (MyErrorHandler) ServeError(w http.ResponseWriter, r *http.Request, err error) {
	weixinLog.Error("Weixin msg error", err.Error())
}
