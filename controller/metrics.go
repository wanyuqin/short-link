package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsController struct {
}

func NewMetricsController() *MetricsController {
	return &MetricsController{}
}

func (ctl *MetricsController) Metrics(c *gin.Context) {
	promhttp.Handler()
}
