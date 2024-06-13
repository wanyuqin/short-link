package controller

import (
	"github.com/gin-gonic/gin"
	"short-link/api/request"
	"short-link/internal/user/services"
)

type UserController struct {
	Controller
}

func NewUserController() *UserController {
	return &UserController{}
}

func (ctl *UserController) Register(c *gin.Context) {
	var req request.Register
	err := c.ShouldBindJSON(&req)
	if err != nil {
		ctl.Error(c, err)
		return
	}
	services.NewUserService().Register(c.Request.Context(), &req)
}

func (ctl *UserController) Login(c *gin.Context) {

}
