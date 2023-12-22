package api

import (
	"delivery_system/client_service/config"
	"delivery_system/pkg/auth"
	"delivery_system/pkg/common_models"
	fasthttputils "delivery_system/pkg/fasthttp_utils"
	"delivery_system/pkg/validator_utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/fasthttp/router"
	"github.com/go-playground/validator/v10"
	"github.com/valyala/fasthttp"
)

const (
	UserIDParam   = "user_id"
	MaskParam     = "mask"
	UsernameParam = "username"
)

type Service interface {
	Register(req common_models.RegisterUserRequest) error
	GetUser(user_id common_models.UserID) (common_models.User, error)
	UpdateUser(req common_models.UpdateUserRequest) (bool, error)
	SearchMask(mask string) ([]common_models.User, error)
	SearchUsername(username string) (common_models.User, bool, error)
	CheckUser(username, password string) (common_models.UserID, bool)
}

type Router struct {
	*router.Router
	cfg *config.Api
	v   *validator.Validate
	s   Service
	a   *auth.Auth
}

func New(cfg *config.Api, service Service) *Router {
	r := Router{
		Router: router.New(),
		cfg:    cfg,
		v:      validator_utils.New(),
		s:      service,
		a:      auth.New(service.CheckUser),
	}

	v1 := r.Group("/api/v1")
	v1.POST("/auth/register", r.registerUser)
	v1.POST("/auth/check", r.checkUser)
	v1.GET(fmt.Sprintf("/users/{%s}", UserIDParam), r.a.AuthMiddleware(r.getUser))
	v1.PATCH(fmt.Sprintf("/users/{%s}", UserIDParam), r.a.AuthMiddleware(r.updateUser))
	v1.GET("/search", r.a.AuthMiddleware(r.search))

	return &r
}

func (r *Router) registerUser(ctx *fasthttp.RequestCtx) {
	// Получение данных
	var req common_models.RegisterUserRequest
	body := ctx.Request.Body()
	if err := json.Unmarshal(body, &req); err != nil {
		fasthttputils.WriteJson(ctx, http.StatusBadRequest, common_models.HttpError{Error: err.Error()})
		return
	}

	// Валидация данных
	if err := r.v.Struct(req); err != nil {
		fasthttputils.WriteJson(ctx, http.StatusBadRequest, common_models.HttpError{Error: err.Error()})
		return
	}

	if err := r.s.Register(req); err != nil {
		fasthttputils.WriteJson(ctx, http.StatusInternalServerError, common_models.HttpError{Error: err.Error()})
		return
	}

	ctx.SetStatusCode(http.StatusOK)
}

func (r *Router) checkUser(ctx *fasthttp.RequestCtx) {
	// Получение данных
	var req auth.CheckAuthRequest
	body := ctx.Request.Body()
	if err := json.Unmarshal(body, &req); err != nil {
		fasthttputils.WriteJson(ctx, http.StatusBadRequest, common_models.HttpError{Error: err.Error()})
		return
	}

	// Валидация данных
	if err := r.v.Struct(req); err != nil {
		fasthttputils.WriteJson(ctx, http.StatusBadRequest, common_models.HttpError{Error: err.Error()})
		return
	}

	userID, ok := r.s.CheckUser(req.Username, req.Password)
	var res auth.CheckAuthResponce
	res.Status = ok
	if ok {
		res.UserID = userID
	}
	fasthttputils.WriteJson(ctx, http.StatusOK, res)
}

func (r *Router) getUser(ctx *fasthttp.RequestCtx) {
	user_id, err := getUserIDParam(ctx)
	if err != nil {
		fasthttputils.WriteJson(ctx, http.StatusBadRequest, common_models.HttpError{Error: err.Error()})
		return
	}

	user, err := r.s.GetUser(user_id)
	if err != nil {
		fasthttputils.WriteJson(ctx, http.StatusInternalServerError, common_models.HttpError{Error: err.Error()})
		return
	}

	fasthttputils.WriteJson(ctx, http.StatusOK, user)
}

func (r *Router) updateUser(ctx *fasthttp.RequestCtx) {
	// Получение данных
	user_id, err := getUserIDParam(ctx)
	if err != nil {
		fasthttputils.WriteJson(ctx, http.StatusBadRequest, common_models.HttpError{Error: err.Error()})
		return
	}
	auth_user_id, err := r.a.GetAuthUserIDValue(ctx)
	if err != nil {
		fasthttputils.WriteJson(ctx, http.StatusBadRequest, common_models.HttpError{Error: err.Error()})
		return
	}
	// Проверяем, что пользователь хочет изменить свои параметры
	if user_id != auth_user_id {
		fasthttputils.WriteJson(ctx, http.StatusForbidden, common_models.HttpError{Error: "auth_user_id != user_id"})
		return
	}

	var req common_models.UpdateUserRequest
	body := ctx.Request.Body()
	if err := json.Unmarshal(body, &req); err != nil {
		fasthttputils.WriteJson(ctx, http.StatusBadRequest, common_models.HttpError{Error: err.Error()})
		return
	}
	req.UserID = user_id

	// Валидация данных
	if err := r.v.Struct(req); err != nil {
		fasthttputils.WriteJson(ctx, http.StatusBadRequest, common_models.HttpError{Error: err.Error()})
		return
	}

	if ok, err := r.s.UpdateUser(req); err != nil {
		fasthttputils.WriteJson(ctx, http.StatusInternalServerError, common_models.HttpError{Error: err.Error()})
		return
	} else if !ok {
		fasthttputils.WriteJson(ctx, http.StatusBadRequest, common_models.HttpError{Error: "such username exists"})
		return
	}

	ctx.SetStatusCode(http.StatusOK)
}

func (r *Router) search(ctx *fasthttp.RequestCtx) {
	// Получение данных
	args := ctx.QueryArgs()
	mask := string(args.Peek(MaskParam))
	username := string(args.Peek(UsernameParam))

	// Валидация данных
	{
		empty := 0
		if len(mask) == 0 {
			empty++
		}
		if len(username) == 0 {
			empty++
		}
		if empty != 1 {
			fasthttputils.WriteJson(ctx, http.StatusBadRequest, common_models.HttpError{Error: "<mask>|<username>"})
			return
		}
	}

	// Поиск пользователей
	var (
		err   error
		users []common_models.User
	)
	if len(mask) != 0 {
		users, err = r.s.SearchMask(mask)
	} else {
		var user common_models.User
		var ok bool
		user, ok, err = r.s.SearchUsername(username)
		if ok {
			users = append(users, user)
		}
	}
	if err != nil {
		fasthttputils.WriteJson(ctx, http.StatusInternalServerError, common_models.HttpError{Error: err.Error()})
		return
	}

	fasthttputils.WriteJson(ctx, http.StatusOK, users)
}

// Возвращает user_id из URL пути .../{user_id}
func getUserIDParam(ctx *fasthttp.RequestCtx) (common_models.UserID, error) {
	user_id, err := strconv.ParseUint(ctx.UserValue(UserIDParam).(string), 10, 64)
	if err != nil {
		return 0, err
	}
	return common_models.UserID(user_id), nil
}
