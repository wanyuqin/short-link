package metrics

import (
	"short-link/internal/consts"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	shortUrlRequestCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: consts.MetricsShortURLRequest,
		Help: "Total number of HTTP requests to the shortUrl endpoint. This metric counts every request made to the shortUrl endpoint, providing insights into the usage frequency of the shortUrl service.",
	}, []string{"shortUrl", "IP"})
)

type ShortUrlRequest struct {
	ShortUrl string `json:"shortUrl"`
	IP       string `json:"IP"`
}

func (sr *ShortUrlRequest) MetricsLabel() prometheus.Labels {
	return prometheus.Labels{
		"shortUrl": sr.ShortUrl,
		"IP":       sr.IP,
	}
}

// RecordShortUrlRequest 记录短链访问指标
func RecordShortUrlRequest(sr *ShortUrlRequest) {
	shortUrlRequestCounter.With(sr.MetricsLabel()).Inc()
}
