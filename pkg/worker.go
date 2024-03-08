package pkg

import (
	"context"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func NewWorker(nc *nats.Conn, js jetstream.JetStream, interval time.Duration) (*Worker, error) {
	streams := map[string]jetstream.Stream{}
	sil := js.ListStreams(context.Background())
	for si := range sil.Info() {
		s, err := js.Stream(context.Background(), si.Config.Name)
		if err != nil {
			return nil, err
		}

		streams[si.Config.Name] = s
	}

	return &Worker{
		nc:       nc,
		js:       js,
		interval: interval,
		streams:  streams,
		metrics:  newStreamMetrics(),
	}, nil
}

type Worker struct {
	nc *nats.Conn
	js jetstream.JetStream

	interval time.Duration

	streams map[string]jetstream.Stream
	metrics *streamMetrics
}

func (w *Worker) Run() error {
	go func() {
		for {
			w.loop()
			time.Sleep(w.interval)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(":2112", nil)
}

func (w *Worker) loop() {
	ctx := context.Background()

	for sn, s := range w.streams {
		if err := w.metrics.collect(ctx, s); err != nil {
			log.Error().Err(err).Msgf("error collecting metrics for stream %s", sn)
		}
	}
}
