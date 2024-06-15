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

// @BasePath /v1/admin
// PingExample godoc
// @Summary 用户注册
// @Schemes
// @Description 用户注册
// @Tags 用户
// @Accept json
// @Param req      body    request.RegisterReq     true        "注册请求参数"
// @Success 200 {object} controller.Response
// @Router /users/register [post]
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

// @BasePath /v1/admin
// PingExample godoc
// @Summary 用户登陆
// @Schemes
// @Description 用户登陆
// @Tags 用户
// @Accept json
// @Param req      body    request.LoginReq     true        "注册请求参数"
// @Success 200 {object} controller.Response
// @Router /users/login [post]
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
