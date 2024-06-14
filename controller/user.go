package controller

import (
	"github.com/gin-gonic/gin"
	"short-link/api/admin/request"
	"short-link/internal/user/services"
)

type UserController struct {
	Controller
}

func NewUserController() *UserController {
	return &UserController{}
}

func (ctl *UserController) Register(c *gin.Context) {
	var req request.RegisterReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	err = services.NewUserService().Register(c.Request.Context(), &req)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Response(c, gin.H{})
}

func (ctl *UserController) Login(c *gin.Context) {
	var req request.LoginReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	err = services.NewUserService().Login(c.Request.Context(), &req)
	if err != nil {

	}
}
