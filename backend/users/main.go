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
	// cfg.host = "0.0.0.0:8080"

	db, err := db.InitDB(cfg.dbaddr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	api := NewUserApi(db)
	srv := InitServer(api)

	if err := srv.Run(cfg.host); err != nil {
		panic(err)
	}
}

func InitServer(api *userApi) *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		v1.POST("/auth/register", api.register)
		v1.GET("/users/:user_id", api.getUser)
		v1.PATCH("/users/:user_id", api.updateUser)
		v1.GET("/search", api.searchUser)
	}

	return r
}

type userApi struct {
	db *db.DB
}

func NewUserApi(db *db.DB) *userApi {
	return &userApi{
		db: db,
	}
}

func (h *userApi) register(c *gin.Context) {
	var req models.RegisterUserRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.Request.Body.Close()

	if err := h.db.AddUser(req); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *userApi) getUser(c *gin.Context) {
	uid, err := getUserIDParam(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.db.GetUser(uid)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, user.User)
}

func (h *userApi) updateUser(c *gin.Context) {
	uid, err := getUserIDParam(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	var req models.UpdateUserRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.Request.Body.Close()

	req.UserID = uid

	if err := h.db.UpdateUser(req); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func getUserIDParam(c *gin.Context) (models.UserID, error) {
	suid := c.Param("user_id")
	if suid == "" {
		return 0, errors.New("empty user_id")
	}
	uid, err := strconv.ParseInt(suid, 10, 64)
	if err != nil {
		return 0, err
	}
	return models.UserID(uid), nil
}

func (h *userApi) searchUser(c *gin.Context) {
	fristname := c.Query("first_name")
	lastname := c.Query("last_name")
	username := c.Query("username")

	users, err := h.db.SearchUser(fristname, lastname, username)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, users)
}
