package admin

import (
	"short-link/controller"

	"github.com/gin-gonic/gin"
)

func NewUserRouter(rg *gin.RouterGroup) {
	usrController := controller.NewUserController()

	userGroup := rg.Group("/users")
	userGroup.POST("/register", usrController.Register)
	userGroup.POST("/login", usrController.Login)
	userGroup.GET("/current-user", usrController.CurrentUser)

}
