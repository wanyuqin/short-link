package app

import (
	"short-link/controller"

	"github.com/gin-gonic/gin"
)

func NewRouter(rg *gin.RouterGroup) {
	linkController := controller.NewLinkController()
	rg.GET("/:short-link", linkController.Request)
}
