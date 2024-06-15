package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"short-link/api/admin/request"
	"short-link/internal/link/services"
)

type LinkController struct {
	Controller
}

func NewLinkController() *LinkController {
	return &LinkController{}
}

// @BasePath /v1/admin
// PingExample godoc
// @Summary 添加链接
// @Schemes
// @Description 添加链接
// @Tags 链接
// @Accept json
// @Param req      body    request.AddLinkReq     true        "注册请求参数"
// @Success 200 {object} controller.Response
// @Router /link/add [post]
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
	shortLink := c.Param("short-link")
	if shortLink == "" {
		c.JSON(http.StatusNotFound, gin.H{})
	}
	url, err := services.NewLinkService().Request(c.Request.Context(), shortLink)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Redirect(c, url)
}
