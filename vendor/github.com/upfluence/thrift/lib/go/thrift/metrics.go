package thrift

import (
	"log"
	"os"
	"time"

	"github.com/upfluence/goutils/tracing"
	"github.com/upfluence/goutils/tracing/noop"
	"github.com/upfluence/goutils/tracing/statsd"
)

var (
	Metrics *Metric = NewMetric(os.Getenv("STATSD_URL"))
)

type Metric struct {
	tracer tracing.Tracer
}

func NewMetric(statsdURL string) *Metric {
	if statsdURL != "" {
		if t, err := statsd.NewTracer(statsdURL, ""); err != nil {
			log.Println("statsd dial: %s", err.Error())
		} else {
			return &Metric{t}
		}
	}
	return &Metric{&noop.Tracer{}}
}

func (m *Metric) Incr(metricName string) {
	m.tracer.Count(metricName, 1)
}

func (m *Metric) Timing(metricName string, duration int64) {
	m.tracer.Timing(metricName, time.Duration(duration))
}
