package auth

import (
	"delivery_system/pkg/common_models"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
)

const AuthUserIDParam = "auth_user_id"

type CheckUserFunc func(username, password string) (common_models.UserID, bool)

type Auth struct {
	checkUser CheckUserFunc
}

func New(checkUser CheckUserFunc) *Auth {
	return &Auth{
		checkUser: checkUser,
	}
}

func NewExternal(host string) *Auth {
	timeout := 2 * time.Second
	maxIdleConnDuration := 24 * time.Hour
	client := &fasthttp.Client{
		ReadTimeout:                   timeout,
		WriteTimeout:                  timeout,
		MaxIdleConnDuration:           maxIdleConnDuration,
		NoDefaultUserAgentHeader:      true, // Don't send: User-Agent: fasthttp
		DisableHeaderNamesNormalizing: true, // If you set the case on your headers correctly you can enable this
		DisablePathNormalizing:        true,
		// increase DNS cache time to an hour instead of default minute
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      1024,
			DNSCacheDuration: time.Hour,
		}).Dial,
	}

	return &Auth{
		checkUser: func(username, password string) (common_models.UserID, bool) {
			reqEntity := &CheckAuthRequest{
				Username: username,
				Password: password,
			}
			reqEntityBytes, _ := json.Marshal(reqEntity)

			req := fasthttp.AcquireRequest()
			req.SetRequestURI(host)
			req.Header.SetMethod(fasthttp.MethodPost)
			req.Header.SetContentTypeBytes([]byte("application/json"))
			req.SetBodyRaw(reqEntityBytes)

			resp := fasthttp.AcquireResponse()
			err := client.Do(req, resp)
			fasthttp.ReleaseRequest(req)
			defer fasthttp.ReleaseResponse(resp)

			if err != nil {
				return 0, false
			}

			statusCode := resp.StatusCode()
			respBody := resp.Body()

			if statusCode != http.StatusOK {
				return 0, false
			}

			respEntity := &CheckAuthResponce{}
			err = json.Unmarshal(respBody, respEntity)
			if err != nil {
				return 0, false
			}

			if !respEntity.Status {
				return 0, false
			}
			return respEntity.UserID, true
		},
	}
}

// func

func (a *Auth) AuthMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		u, p, ok := basicAuth(ctx)
		if ok {
			if userID, ok := a.checkUser(u, p); ok {
				ctx.SetUserValue(AuthUserIDParam, userID)
				next(ctx)
				return
			}
		}

		// Request Basic Authentication otherwise
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized), fasthttp.StatusUnauthorized)
		ctx.Response.Header.Set("WWW-Authenticate", "Basic realm=Restricted")
	}
}

func (a *Auth) GetAuthUserIDValue(ctx *fasthttp.RequestCtx) (common_models.UserID, error) {
	user_id, ok := ctx.UserValue(AuthUserIDParam).(common_models.UserID)
	if !ok {
		return 0, errors.New("auth user_id is not uint64")
	}
	return common_models.UserID(user_id), nil
}
