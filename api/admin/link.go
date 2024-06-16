package admin

import (
	"github.com/gin-gonic/gin"
	"short-link/controller"
)

func NewLinkRouter(rg *gin.RouterGroup) {
	linkController := controller.NewLinkController()

	linkGroup := rg.Group("/link")
	{
		linkGroup.POST("/add", linkController.AddLink)
		linkGroup.POST("/list", linkController.LinkList)
	}
}
