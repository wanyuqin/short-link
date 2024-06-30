package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"short-link/api/admin/request"
	"short-link/ctxkit"
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
	if err := req.Validate(); err != nil {
		ctl.ParamException(c, err)
		return
	}
	userId := ctxkit.GetUserId(c.Request.Context())
	if userId == 0 {
		ctl.UnauthorizedException(c)
		return
	}
	req.UserId = userId
	err := services.NewLinkService().AddLink(c.Request.Context(), &req)
	if err != nil {
		ctl.Error(c, err)
		return
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
		ctl.PageNotFound(c)
		return
	}
	ctl.Redirect(c, url)
}

func (ctl *LinkController) LinkList(c *gin.Context) {
	var req request.LinkListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ctl.ParamException(c, err)
		return
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	resp, err := services.NewLinkService().LinkList(c.Request.Context(), &req)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Response(c, resp)
}

func (ctl *LinkController) DeleteLink(c *gin.Context) {
	var req request.DelLinkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ctl.ParamException(c, err)
		return
	}
	if err := req.Validate(); err != nil {
		ctl.ParamException(c, err)
		return
	}
	userId := ctxkit.GetUserId(c.Request.Context())
	if userId == 0 {
		ctl.UnauthorizedException(c)
		return
	}
	req.UserId = userId
	err := services.NewLinkService().DeleteLink(c.Request.Context(), &req)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Response(c, gin.H{})
}
