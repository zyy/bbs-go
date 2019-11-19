package weixin

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type MyErrorHandler struct{}

func (MyErrorHandler) ServeError(w http.ResponseWriter, r *http.Request, err error) {
	logrus.Error("Weixin msg error", err.Error())
}
