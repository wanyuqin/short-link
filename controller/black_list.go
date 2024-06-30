package controller

import (
	"github.com/gin-gonic/gin"
	"short-link/api/admin/request"
	"short-link/ctxkit"
	"short-link/internal/link/services"
)

type BlackListController struct {
	Controller
}

func NewBlackListController() *BlackListController {
	return &BlackListController{}
}

func (ctl *BlackListController) AddBlackList(c *gin.Context) {
	var req request.AddBlackListReq
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
	if err := services.NewBlackListService().AddBlackList(c.Request.Context(), &req); err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Response(c, gin.H{})
	return
}

func (ctl *BlackListController) DeleteBlackList(c *gin.Context) {
	var req request.DeleteBlackListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ctl.ParamException(c, err)
		return
	}
	if err := req.Validate(); err != nil {
		ctl.ParamException(c, err)
		return
	}
	if err := services.NewBlackListService().DeleteBlackList(c.Request.Context(), req.Id); err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Response(c, gin.H{})
	return
}

func (ctl *BlackListController) ListBlackList(c *gin.Context) {
	var req request.ListBlackListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ctl.ParamException(c, err)
		return
	}
	if err := req.Validate(); err != nil {
		ctl.ParamException(c, err)
		return
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	resp, err := services.NewBlackListService().ListBlackList(c.Request.Context(), &req)
	if err != nil {
		ctl.Error(c, err)
		return
	}

	ctl.Response(c, resp)
	return
}
