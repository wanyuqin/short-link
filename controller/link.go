package controller

import (
	"github.com/gin-gonic/gin"
	"short-link/api/request"
)

type LinkController struct {
	Controller
}

func NewLinkController() *LinkController {
	return &LinkController{}
}

func (ctl *LinkController) AddLink(c *gin.Context) {
	var req request.AddLinkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ctl.Error(c, err)
		return
	}

}
