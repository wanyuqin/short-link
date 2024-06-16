package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Controller struct {
}

type Response struct {
	Code      int         `json:"code"`
	Data      interface{} `json:"data"`
	Msg       string      `json:"msg"`
	TimeStamp int64       `json:"timeStamp"`
}

func (ctl *Controller) Response(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:      http.StatusOK,
		Data:      data,
		Msg:       "success",
		TimeStamp: time.Now().UnixMilli(),
	})
}

func (ctl *Controller) Error(c *gin.Context, err error) {
	c.JSON(http.StatusOK, Response{
		Code:      http.StatusInternalServerError,
		Msg:       err.Error(),
		TimeStamp: time.Now().Unix(),
	})
	c.Abort()
}

func (ctl *Controller) Redirect(c *gin.Context, url string) {
	if url == "" {
		c.JSON(http.StatusNotFound, gin.H{})
		c.Abort()
	}
	c.Redirect(http.StatusFound, url)
	c.Abort()
}

type UserToken struct {
	Id       uint64 `json:"id"`
	Username string `json:"username"`
}
