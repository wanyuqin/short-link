package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"short-link/internal/consts"
)

var ()

type ShortUrlRequest struct {
	ShortUrl string `json:"shortUrl"`
	Ip       string `json:"ip"`
}

func (sr *ShortUrlRequest) MetricsLabel() prometheus.Labels {
	return prometheus.Labels{
		"shortUrl": sr.ShortUrl,
		"ip":       sr.Ip,
	}
}

func RecordShortUrlRequest(sr *ShortUrlRequest) {
	shortUrlRequestCounter := promauto.NewCounter(prometheus.CounterOpts{
		Name:        consts.MetricsShortUrlRequest,
		Help:        "Total number of HTTP requests to the shortUrl endpoint. This metric counts every request made to the shortUrl endpoint, providing insights into the usage frequency of the shortUrl service.",
		ConstLabels: sr.MetricsLabel(),
	},
	)
	shortUrlRequestCounter.Inc()
}
