package pkg

import (
	"context"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type streamMetrics struct {
	streamBytes     *prometheus.GaugeVec
	streamMsgs      *prometheus.GaugeVec
	streamConsumers *prometheus.GaugeVec
	streamSubjects  *prometheus.GaugeVec

	consumerAckPending  *prometheus.GaugeVec
	consumerRedelivered *prometheus.GaugeVec
	consumerWaiting     *prometheus.GaugeVec
	consumerPending     *prometheus.GaugeVec
}

func newStreamMetrics() *streamMetrics {
	streamLabels := []string{"stream"}
	consumerLabels := []string{"stream", "consumer"}

	return &streamMetrics{
		streamBytes: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "nats_stream_bytes",
			Help: "the number of bytes stored in the stream",
		}, streamLabels),
		streamMsgs: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "nats_stream_messages",
			Help: "the number of messages stored in the stream",
		}, streamLabels),
		streamConsumers: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "nats_stream_consumers",
			Help: "the number of consumers on the stream",
		}, streamLabels),
		streamSubjects: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "nats_stream_subjects",
			Help: "the number of unique subjects the stream has received messages on",
		}, streamLabels),

		consumerAckPending: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "nats_consumer_ack_pending",
			Help: "the number of messages that have been delivered but not yet acknowledged",
		}, consumerLabels),
		consumerRedelivered: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "nats_consumer_redelivered",
			Help: "the number of messages that have been redelivered and not yet acknowledged. Each message is counted only once, even if it has been redelivered multiple times. This count is reset when the message is eventually acknowledged.",
		}, consumerLabels),
		consumerWaiting: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "nats_consumer_waiting",
			Help: "the count of active pull requests. It is only relevant for pull-based consumers",
		}, consumerLabels),
		consumerPending: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "nats_consumer_pending",
			Help: "the number of messages that match the consumer's filter, but have not been delivered yet",
		}, consumerLabels),
	}
}

func (m *streamMetrics) collect(ctx context.Context, s jetstream.Stream) error {
	si, err := s.Info(ctx)
	if err != nil {
		return err
	}

	m.streamBytes.WithLabelValues(si.Config.Name).Set(float64(si.State.Bytes))
	m.streamMsgs.WithLabelValues(si.Config.Name).Set(float64(si.State.Msgs))
	m.streamConsumers.WithLabelValues(si.Config.Name).Set(float64(si.State.Consumers))
	m.streamSubjects.WithLabelValues(si.Config.Name).Set(float64(si.State.NumSubjects))

	consumers := s.ListConsumers(ctx)
	for cons := range consumers.Info() {
		m.consumerAckPending.WithLabelValues(si.Config.Name, cons.Name).Set(float64(cons.NumAckPending))
		m.consumerRedelivered.WithLabelValues(si.Config.Name, cons.Name).Set(float64(cons.NumRedelivered))
		m.consumerWaiting.WithLabelValues(si.Config.Name, cons.Name).Set(float64(cons.NumWaiting))
		m.consumerPending.WithLabelValues(si.Config.Name, cons.Name).Set(float64(cons.NumPending))
	}

	return nil
}
