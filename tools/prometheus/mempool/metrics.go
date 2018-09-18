package mempool

import (
	"github.com/irisnet/irishub/client/context"
	"github.com/tendermint/tendermint/mempool"
	"log"
	"time"
)

// Metrics contains metrics exposed by this package.
// see MetricsProvider for descriptions.
type Metrics struct {
	TmMetrics mempool.Metrics
}

// PrometheusMetrics returns Metrics build using Prometheus client library.
func PrometheusMetrics() *Metrics {
	tmMetrics := *mempool.PrometheusMetrics()
	return &Metrics{
		tmMetrics,
	}
}

func (m *Metrics) Start(rpc context.CLIContext) {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			if result, err := rpc.NumUnconfirmedTxs(); err == nil {
				m.TmMetrics.Size.Set(float64(result.N))
			} else {
				log.Println(err)
			}
		}
	}()
}
