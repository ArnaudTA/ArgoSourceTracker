package metrics

import (
	"time"

	"github.com/cableship/chart-sentinel/pkg/argocd/application"
	"github.com/cableship/chart-sentinel/pkg/config"
)

func StartAppMonitor() {

	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				generateMetrics()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	generateMetrics()

}

func generateMetrics() {
	sum := application.List([]string{}, "", 0, 10000)
	stats := sum.Stats
	for status, n := range stats {
		Metrics.GetMetric(AppsStatusCountMetricName).SetGaugeValue([]string{string(status), config.Global.Argocd.Instance}, float64(n))
	}
}
