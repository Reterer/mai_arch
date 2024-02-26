package main

import (
	"delivery/common/db"
	"delivery/common/models"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type config struct {
	dbaddr string
	host   string
}

func getConfig() *config {
	const (
		DB_ADDR = "DB_ADDR"
		HOST    = "HOST"
	)

	return &config{
		dbaddr: os.Getenv(DB_ADDR),
		host:   os.Getenv(HOST),
	}
}

func main() {
	cfg := getConfig()

	// cfg.dbaddr = "root:password@tcp(localhost:3306)/delivery"
	// cfg.host = "localhost:8083"

	db, err := db.InitDB(cfg.dbaddr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	api := NewDeliveryApi(db)
	srv := InitServer(api)

	if err := srv.Run(cfg.host); err != nil {
		panic(err)
	}
}

func InitServer(api *deliveryApi) *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		v1.POST("/deliveries", api.crateDelivery)
		v1.GET("/deliveries_from/:user_id", api.deliveriesFrom)
		v1.GET("/deliveries_to/:user_id", api.deliveriesTo)
	}

	return r
}

type deliveryApi struct {
	db *db.DB
}

func NewDeliveryApi(db *db.DB) *deliveryApi {
	return &deliveryApi{
		db: db,
	}
}

func (h *deliveryApi) crateDelivery(c *gin.Context) {
	var req models.CreateDeliveryRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.Request.Body.Close()

	if err := h.db.AddDelivery(req); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *deliveryApi) deliveriesFrom(c *gin.Context) {
	uid, err := getUserIDParam(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	deliveries, err := h.db.GetDeliveriesByFromUserID(uid)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, deliveries)
}

func (h *deliveryApi) deliveriesTo(c *gin.Context) {
	uid, err := getUserIDParam(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	deliveries, err := h.db.GetDeliveriesByToUserID(uid)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, deliveries)
}

func getUserIDParam(c *gin.Context) (models.UserID, error) {
	uid, err := getIntParam(c, "user_id")
	return models.UserID(uid), err
}

func getIntParam(c *gin.Context, param string) (int64, error) {
	sparam := c.Param(param)
	if sparam == "" {
		return 0, errors.New("empty " + param)
	}
	iparam, err := strconv.ParseInt(sparam, 10, 64)
	if err != nil {
		return 0, err
	}
	return iparam, nil
}
