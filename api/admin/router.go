package admin

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(engine *gin.RouterGroup) {
	v1 := engine.Group("/api/v1/admin")
	NewUserRouter(v1)
	NewLinkRouter(v1)
}
