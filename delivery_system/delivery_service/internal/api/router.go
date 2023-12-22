package api

import (
	"delivery_system/delivery_service/config"
	"delivery_system/pkg/auth"
	"delivery_system/pkg/validator_utils"

	"github.com/fasthttp/router"
	"github.com/go-playground/validator/v10"
)

const (
	UserIDParam   = "user_id"
	MaskParam     = "mask"
	UsernameParam = "username"
)

type Service interface {
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
		// a:      auth.New(service.CheckUser), // TODO
	}

	// v1 := r.Group("/api/v1")
	// v1.POST("/auth/register", r.registerUser)
	// v1.POST("/auth/check", r.checkUser)
	// v1.GET(fmt.Sprintf("/users/{%s}", UserIDParam), r.a.AuthMiddleware(r.getUser))
	// v1.PATCH(fmt.Sprintf("/users/{%s}", UserIDParam), r.a.AuthMiddleware(r.updateUser))
	// v1.GET("/search", r.a.AuthMiddleware(r.search))

	return &r
}
