package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
}

func (ctl *Controller) Response() {

}

func (ctl *Controller) Error(c *gin.Context, err error) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusInternalServerError,
		"data": "",
		"msg":  err.Error(),
	})
	c.Abort()
}
