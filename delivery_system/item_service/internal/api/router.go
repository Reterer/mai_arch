package api

import (
	"delivery_system/item_service/config"
	"delivery_system/pkg/auth"
	"delivery_system/pkg/common_models"
	fasthttputils "delivery_system/pkg/fasthttp_utils"
	"delivery_system/pkg/validator_utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/fasthttp/router"
	"github.com/go-playground/validator/v10"
	"github.com/valyala/fasthttp"
)

const (
	ItemIDParam   = "item_id"
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
	a   *auth.Auth
}

func New(cfg *config.Api, service Service) *Router {
	r := Router{
		Router: router.New(),
		cfg:    cfg,
		v:      validator_utils.New(),
		s:      service,
		a:      auth.New(auth.NewDefaultChecker()), // TODO
	}

	v1 := r.Group("/api/v1")
	v1.POST("/items", r.a.AuthMiddleware(r.createItem)) // Нужна авторизация
	v1.GET("/items/{item_id}", r.getItem)
	v1.PATCH("/items/{item_id}", r.a.AuthMiddleware(r.patchItem)) // Нужна авторизация
	v1.GET("/items_by_username/{username}", r.getItemsByUsername)

	return &r
}

func (r *Router) createItem(ctx *fasthttp.RequestCtx) {
	userID, err := r.a.GetAuthUserIDValue(ctx)
	if err != nil {
		fasthttputils.WriteJson(ctx, http.StatusUnauthorized, common_models.HttpError{Error: err.Error()})
		return
	}
	// Получение данных
	var req common_models.CreateItemRequest
	body := ctx.Request.Body()
	if err := json.Unmarshal(body, &req); err != nil {
		fasthttputils.WriteJson(ctx, http.StatusBadRequest, common_models.HttpError{Error: err.Error()})
		return
	}
	req.OwnerID = userID

	// Валидация данных
	if err := r.v.Struct(req); err != nil {
		fasthttputils.WriteJson(ctx, http.StatusBadRequest, common_models.HttpError{Error: err.Error()})
		return
	}

	// Выполнение запроса
	err = r.s.CreateItem(req)
	if err != nil {
		fasthttputils.WriteJson(ctx, http.StatusInternalServerError, common_models.HttpError{Error: err.Error()})
		return
	}
	ctx.SetStatusCode(http.StatusOK)
}

func (r *Router) getItem(ctx *fasthttp.RequestCtx) {
	itemID, err := getItemIDParam(ctx)
	if err != nil {
		fasthttputils.WriteJson(ctx, http.StatusBadRequest, common_models.HttpError{Error: err.Error()})
		return
	}
	// Выполнение запроса
	item, err := r.s.GetItem(itemID)
	if err != nil {
		fasthttputils.WriteJson(ctx, http.StatusInternalServerError, common_models.HttpError{Error: err.Error()})
		return
	}

	fasthttputils.WriteJson(ctx, http.StatusOK, item)
}

func (r *Router) patchItem(ctx *fasthttp.RequestCtx) {
	userID, err := r.a.GetAuthUserIDValue(ctx)
	if err != nil {
		fasthttputils.WriteJson(ctx, http.StatusUnauthorized, common_models.HttpError{Error: err.Error()})
		return
	}

	itemID, err := getItemIDParam(ctx)
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
		ItemID:  itemID,
		Data:    req.Data,
		OwnerID: userID,
	}

	// Выполнение запроса
	err = r.s.UpdateItem(item)
	if err != nil {
		fasthttputils.WriteJson(ctx, http.StatusInternalServerError, common_models.HttpError{Error: err.Error()})
		return
	}

	ctx.SetStatusCode(http.StatusOK)
}

func (r *Router) getItemsByUsername(ctx *fasthttp.RequestCtx) {
	username := ctx.UserValue(UsernameParam).(string)

	// Выполнение запроса
	items, err := r.s.GetItemsByUsername(username)
	if err != nil {
		fasthttputils.WriteJson(ctx, http.StatusInternalServerError, common_models.HttpError{Error: err.Error()})
		return
	}

	fasthttputils.WriteJson(ctx, http.StatusOK, items)
}

// Возвращает item_id из URL пути .../{item_id}
func getItemIDParam(ctx *fasthttp.RequestCtx) (common_models.ItemID, error) {
	itemID, err := strconv.ParseUint(ctx.UserValue(ItemIDParam).(string), 10, 64)
	if err != nil {
		return 0, err
	}
	return common_models.ItemID(itemID), nil
}
