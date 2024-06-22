package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Controller struct {
}

type Response struct {
	Code      int         `json:"code,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Msg       string      `json:"msg,omitempty"`
	TimeStamp int64       `json:"timeStamp,omitempty"`
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

func (ctl *Controller) ParamException(c *gin.Context, err error) {
	c.JSON(http.StatusOK, Response{
		Code:      http.StatusBadRequest,
		Msg:       err.Error(),
		TimeStamp: time.Now().Unix(),
	})
	c.Abort()
}

func (ctl *Controller) UnauthorizedException(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code:      http.StatusUnauthorized,
		Msg:       "用户未登录",
		TimeStamp: time.Now().Unix(),
	})
	c.Abort()
}

func (ctl *Controller) PageNotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, Response{
		Msg: "page not found",
	})
}

type UserToken struct {
	Id       uint64 `json:"id"`
	Username string `json:"username"`
}
