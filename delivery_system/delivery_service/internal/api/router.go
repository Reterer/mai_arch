package api

import (
	"delivery_system/delivery_service/config"
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
	UserIDParam = "user_id"
)

type Service interface {
	CreateDelivery(req common_models.CreateDeliveryRequest) error
	GetDeliveriesByFrom(userID common_models.UserID) ([]common_models.Delivery, error)
	GetDeliveriesByTo(userID common_models.UserID) ([]common_models.Delivery, error)
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
	v1.POST("/deliveries", r.createDelivery)
	v1.GET("/deliveries_from/{user_id}", r.deliveriesByFrom)
	v1.GET("/deliveries_to/{user_id}", r.deliveriesByTo)

	return &r
}

func (r *Router) createDelivery(ctx *fasthttp.RequestCtx) {
	// Получение данных
	var req common_models.CreateDeliveryRequest
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
	if len(req.Items) == 0 {
		fasthttputils.WriteJson(ctx, http.StatusBadRequest, common_models.HttpError{Error: "empty items"})
		return
	}

	// Выполнение запроса
	err := r.s.CreateDelivery(req)
	if err != nil {
		fasthttputils.WriteJson(ctx, http.StatusInternalServerError, common_models.HttpError{Error: err.Error()})
		return
	}
	ctx.SetStatusCode(http.StatusOK)
}

func (r *Router) deliveriesByFrom(ctx *fasthttp.RequestCtx) {
	userID, err := fasthttputils.GetUint64Param(ctx, UserIDParam)
	if err != nil {
		fasthttputils.WriteJson(ctx, http.StatusBadRequest, common_models.HttpError{Error: err.Error()})
		return
	}

	// Выполнение запроса
	deliveries, err := r.s.GetDeliveriesByFrom(common_models.UserID(userID))
	if err != nil {
		fasthttputils.WriteJson(ctx, http.StatusInternalServerError, common_models.HttpError{Error: err.Error()})
		return
	}
	fasthttputils.WriteJson(ctx, http.StatusOK, deliveries)
}

func (r *Router) deliveriesByTo(ctx *fasthttp.RequestCtx) {
	userID, err := fasthttputils.GetUint64Param(ctx, UserIDParam)
	if err != nil {
		fasthttputils.WriteJson(ctx, http.StatusBadRequest, common_models.HttpError{Error: err.Error()})
		return
	}

	// Выполнение запроса
	deliveries, err := r.s.GetDeliveriesByTo(common_models.UserID(userID))
	if err != nil {
		fasthttputils.WriteJson(ctx, http.StatusInternalServerError, common_models.HttpError{Error: err.Error()})
		return
	}
	fasthttputils.WriteJson(ctx, http.StatusOK, deliveries)
}
