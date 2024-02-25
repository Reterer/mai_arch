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

	cfg.dbaddr = "root:password@tcp(localhost:3306)/delivery"
	cfg.host = "localhost:8081"

	db, err := db.InitDB(cfg.dbaddr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	api := NewItemApi(db)
	srv := InitServer(api)

	if err := srv.Run(cfg.host); err != nil {
		panic(err)
	}
}

func InitServer(api *itemApi) *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		v1.POST("/items", api.createItem)
		v1.GET("/items/:item_id", api.getItem)
		v1.PATCH("/items/:item_id", api.updateItem)
		v1.GET("/items_by_user/:user_id", api.getItemsByUserID)
	}

	return r
}

type itemApi struct {
	db *db.DB
}

func NewItemApi(db *db.DB) *itemApi {
	return &itemApi{
		db: db,
	}
}

func (h *itemApi) createItem(c *gin.Context) {
	var req models.CreateItemRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.Request.Body.Close()

	if err := h.db.AddItem(req); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *itemApi) getItem(c *gin.Context) {
	itemID, err := getItemIDParam(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	item, err := h.db.GetItem(itemID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *itemApi) updateItem(c *gin.Context) {
	itemID, err := getItemIDParam(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	var req models.UpdateItemRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.Request.Body.Close()

	req.ItemID = itemID

	if err := h.db.UpdateItem(req); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *itemApi) getItemsByUserID(c *gin.Context) {
	uid, err := getUserIDParam(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	items, err := h.db.GetItemsByUserID(uid)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, items)
}

func getItemIDParam(c *gin.Context) (models.ItemID, error) {
	itemID, err := getIntParam(c, "item_id")
	return models.ItemID(itemID), err
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
