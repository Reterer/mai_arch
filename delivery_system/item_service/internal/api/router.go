package api

import (
	"delivery_system/item_service/config"
	"delivery_system/pkg/common_models"
	fasthttputils "delivery_system/pkg/fasthttp_utils"
	"delivery_system/pkg/validator_utils"
	"encoding/json"
	"net/http"

	"github.com/fasthttp/router"
	"github.com/go-playground/validator/v10"
	"github.com/valyala/fasthttp"
)

const (
	ItemIDParam   = "item_id"
	UserIDParam   = "user_id"
	UsernameParam = "username"
)

type Service interface {
	CreateItem(newItem common_models.CreateItemRequest) error
	GetItem(itemID common_models.ItemID) (common_models.Item, error)
	UpdateItem(item common_models.Item) error
	GetItemsByUsername(username string) ([]common_models.Item, error)
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
	v1.POST("/items", r.createItem)
	v1.GET("/items/{item_id}", r.getItem)
	v1.PATCH("/items/{item_id}", r.patchItem)
	v1.GET("/items_by_user_id/{user_id}", r.getItemsByUserID)

	return &r
}

func (r *Router) createItem(ctx *fasthttp.RequestCtx) {
	// Получение данных
	var req common_models.CreateItemRequest
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

	// Выполнение запроса
	err := r.s.CreateItem(req)
	if err != nil {
		fasthttputils.WriteJson(ctx, http.StatusInternalServerError, common_models.HttpError{Error: err.Error()})
		return
	}
	ctx.SetStatusCode(http.StatusOK)
}

func (r *Router) getItem(ctx *fasthttp.RequestCtx) {
	itemID, err := fasthttputils.GetUint64Param(ctx, ItemIDParam)
	if err != nil {
		fasthttputils.WriteJson(ctx, http.StatusBadRequest, common_models.HttpError{Error: err.Error()})
		return
	}
	// Выполнение запроса
	item, err := r.s.GetItem(common_models.ItemID(itemID))
	if err != nil {
		fasthttputils.WriteJson(ctx, http.StatusInternalServerError, common_models.HttpError{Error: err.Error()})
		return
	}

	fasthttputils.WriteJson(ctx, http.StatusOK, item)
}

func (r *Router) patchItem(ctx *fasthttp.RequestCtx) {
	itemID, err := fasthttputils.GetUint64Param(ctx, ItemIDParam)
	if err != nil {
		fasthttputils.WriteJson(ctx, http.StatusBadRequest, common_models.HttpError{Error: err.Error()})
		return
	}

	var req common_models.PatchItemRequest
	body := ctx.Request.Body()
	if err := json.Unmarshal(body, &req); err != nil {
		fasthttputils.WriteJson(ctx, http.StatusBadRequest, common_models.HttpError{Error: err.Error()})
		return
	}

	item := common_models.Item{
		ItemID:  common_models.ItemID(itemID),
		Data:    req.Data,
		OwnerID: req.OwnerID,
	}

	// Выполнение запроса
	err = r.s.UpdateItem(item)
	if err != nil {
		fasthttputils.WriteJson(ctx, http.StatusInternalServerError, common_models.HttpError{Error: err.Error()})
		return
	}

	ctx.SetStatusCode(http.StatusOK)
}

func (r *Router) getItemsByUserID(ctx *fasthttp.RequestCtx) {
	username := ctx.UserValue(UsernameParam).(string)

	// Выполнение запроса
	items, err := r.s.GetItemsByUsername(username)
	if err != nil {
		fasthttputils.WriteJson(ctx, http.StatusInternalServerError, common_models.HttpError{Error: err.Error()})
		return
	}

	fasthttputils.WriteJson(ctx, http.StatusOK, items)
}
