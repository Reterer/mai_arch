package fasthttputils

import (
	"strconv"

	"github.com/valyala/fasthttp"
)

// Возвращает uint64 param из URL пути .../{user_id}
func GetUint64Param(ctx *fasthttp.RequestCtx, param string) (uint64, error) {
	res, err := strconv.ParseUint(ctx.UserValue(param).(string), 10, 64)
	if err != nil {
		return 0, err
	}
	return res, nil
}
