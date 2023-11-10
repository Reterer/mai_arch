package fasthttputils

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func WriteJson[T any](ctx *fasthttp.RequestCtx, code int, msg T) {
	errMsg, err := json.Marshal(msg)
	if err != nil {
		zap.L().Error("unmarshal", zap.Error(err))
	}

	ctx.SetStatusCode(code)
	ctx.SetContentType("application/json")
	_, err = ctx.Write(errMsg)
	if err != nil {
		zap.L().Error("write", zap.Error(err))
	}
}
