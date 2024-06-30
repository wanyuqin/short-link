package admin

import (
	"github.com/gin-gonic/gin"
	"short-link/controller"
)

func NewLinkRouter(rg *gin.RouterGroup) {
	linkController := controller.NewLinkController()
	blackListController := controller.NewBlackListController()

	linkGroup := rg.Group("/link")
	{
		linkGroup.POST("/add", linkController.AddLink)
		linkGroup.POST("/list", linkController.LinkList)
		linkGroup.POST("/del", linkController.DeleteLink)
	}

	linkBlackListGroup := linkGroup.Group("/black-list")
	{
		linkBlackListGroup.POST("/add", blackListController.AddBlackList)
		linkBlackListGroup.POST("/del", blackListController.DeleteBlackList)
		linkBlackListGroup.POST("/list", blackListController.ListBlackList)

	}
}
