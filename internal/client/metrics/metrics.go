package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type ClientMetrics struct {
	DelayTime      *prometheus.HistogramVec
	ReceiveCounter *prometheus.CounterVec
}

func NewClientMetrics(namespace string) ClientMetrics {
	var cm ClientMetrics

	cm.DelayTime = promauto.NewHistogramVec(
		prometheus.HistogramOpts{ //nolint:exhaustivestruct
			Namespace: namespace,
			Subsystem: "quic",
			Name:      "quic_delay",
			Help:      "Quic Delay",
		},
		[]string{
			"client",
		},
	)

	cm.ReceiveCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "quic",
			Name:      "quic_receive_msg_count",
			Help:      "quic received message count",
		},
		[]string{
			"client",
		},
	)

	return cm
}

//func (m ClientMetrics) AddDelayTime(sample float64) {
//	m.DelayTime.With(c)
//}

func (m ClientMetrics) IncReceive(client string) {
	m.ReceiveCounter.With(prometheus.Labels{"client": client}).Inc()
}
