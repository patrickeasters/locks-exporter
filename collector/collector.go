package collector

import (
	"io"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

const Namespace = "locks"

var _ prometheus.Collector = (*Collector)(nil)

// Exporter collects metrics from a local Raspberry Pi
type Collector struct {
	logger         *logrus.Logger
	throttleStatus io.Reader

	podFileLocks *prometheus.Desc
}

type lock struct {
	count     int
	pod       string
	container string
	namespace string
}

// New returns an initialized collector
func New(logger *logrus.Logger) *Collector {
	return &Collector{
		logger: logger,
		podFileLocks: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "pod", "file_locks"),
			"Number of file locks held by processes in pod",
			[]string{"pod", "namespace"},
			nil,
		),
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

// Collect fetches the statistics from the configured memcached server, and
// delivers them as Prometheus metrics. It implements prometheus.Collector.
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	locks := getLocks()

	for _, l := range locks {
		ch <- prometheus.MustNewConstMetric(
			c.podFileLocks,
			prometheus.GaugeValue,
			float64(l.count),
			l.pod,
			l.namespace,
		)
	}
}
