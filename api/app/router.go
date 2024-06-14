package app

import (
	"github.com/gin-gonic/gin"
	"short-link/controller"
)

func NewRouter(rg *gin.RouterGroup) {
	linkController := controller.NewLinkController()
	rg.GET("/", linkController.Request)
}
