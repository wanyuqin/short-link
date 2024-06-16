package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"short-link/api/admin/request"
	"short-link/api/middleware"
	"short-link/internal/user/services"
	"strings"
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
	user, err := services.NewUserService().Login(c.Request.Context(), &req)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	token, _ := middleware.GetToken(user.ID, user.Username)

	ctl.Response(c, gin.H{"token": token})
}

func (ctl *UserController) CurrentUser(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		ctl.Error(c, errors.New("未知用户"))
		return
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		ctl.Error(c, errors.New("请求头中auth格式有误"))
		return
	}
	m, err := middleware.ParseToken(parts[1])
	if err != nil {
		ctl.Error(c, err)
		return
	}
	ctl.Response(c, m)
}
