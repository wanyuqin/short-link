package api

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(engine *gin.RouterGroup) {
	v1 := engine.Group("/v1")
	NewUserRouter(v1)
}
