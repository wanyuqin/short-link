package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
}

func (ctl *Controller) Response(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": data,
		"msg":  "success",
	})
}

func (ctl *Controller) Error(c *gin.Context, err error) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusInternalServerError,
		"data": "",
		"msg":  err.Error(),
	})
	c.Abort()
}
