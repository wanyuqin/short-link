package controller

import (
	"github.com/gin-gonic/gin"
	"short-link/api/admin/request"
	"short-link/internal/link/services"
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
	err := services.NewLinkService().AddLink(c.Request.Context(), &req)
	if err != nil {
		ctl.Error(c, err)
	}
	ctl.Response(c, gin.H{})
	return
}

func (ctl *LinkController) Request(c *gin.Context) {

}
