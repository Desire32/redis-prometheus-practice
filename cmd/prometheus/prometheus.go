package prometheus

import "github.com/prometheus/client_golang/prometheus"

var RedisPopTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "redis_pop_total",
		Help: "deleting objects",
	},
	[]string{"objects"},
)

func init() {
	prometheus.MustRegister(RedisPopTotal)
}
