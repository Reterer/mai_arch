package api

import (
	"delivery_system/client_service/config"
	"delivery_system/pkg/common_models"
	fasthttputils "delivery_system/pkg/fasthttp_utils"
	"delivery_system/pkg/validator_utils"
	"encoding/json"
	"net/http"

	"github.com/fasthttp/router"
	"github.com/go-playground/validator/v10"
	"github.com/valyala/fasthttp"
)

type Service interface {
}

type Router struct {
	*router.Router
	cfg *config.Api
	v   *validator.Validate
	s   Service
}

func New(cfg *config.Api, service Service) *Router {
	r := Router{
		Router: router.New(),
		cfg:    cfg,
		v:      validator_utils.New(),
		s:      service,
	}

	v1 := r.Group("/api/v1")
	v1.POST("/auth/register", r.registerUser)
	v1.GET("/users/{user_id}", authMiddleware(r.getUser))
	v1.PATCH("/users/{user_id}", authMiddleware(r.updateUser))
	v1.GET("/search", authMiddleware(r.search))

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

	// TODO Создание пользователя

	fasthttputils.WriteJson(ctx, http.StatusOK, "ok")
}

func (r *Router) getUser(ctx *fasthttp.RequestCtx) {
	// TODO логика получения пользователя
	var user common_models.User

	fasthttputils.WriteJson(ctx, http.StatusOK, user)
}

func (r *Router) updateUser(ctx *fasthttp.RequestCtx) {
	// Получение данных
	var req common_models.UpdateUserRequest
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

	// TODO обновление пользователя

	fasthttputils.WriteJson(ctx, http.StatusOK, "ok")
}

func (r *Router) search(ctx *fasthttp.RequestCtx) {
	// Получение данных
	args := ctx.QueryArgs()
	mask := string(args.Peek("mask"))
	username := string(args.Peek("username"))

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

	// TODO поиск пользователей
	var users []common_models.User

	fasthttputils.WriteJson(ctx, http.StatusOK, users)
}
