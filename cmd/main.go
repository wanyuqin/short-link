package main

import (
	"github.com/gin-gonic/gin"
	"short-link/api"
)

func main() {
	engine := gin.Default()
	api.NewRouter(engine)
	engine.Run(":8080")
}
