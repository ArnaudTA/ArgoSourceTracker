package metrics

import (
	"argocd-watcher/pkg/config"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/sirupsen/logrus"
)

var Metrics *ginmetrics.Monitor
var AppsStatusCountMetricName = "ast_apps_status_count"

func Init() { // get global Monitor object
	Metrics = ginmetrics.GetMonitor()
	Metrics.AddMetric(&ginmetrics.Metric{
		Type:        ginmetrics.Gauge,
		Name:        AppsStatusCountMetricName,
		Description: fmt.Sprintf("Apps status count"),
		Labels: []string{
			"status",
			"instance",
		},
	})
	StartAppMonitor()
}
func Expose(r *gin.Engine) {
	// +optional set metric path, default /debug/metrics
	Metrics.SetMetricPath("/metrics")
	// +optional set slow time, default 5s
	Metrics.SetSlowTime(1)
	// +optional set request duration, default {0.1, 0.3, 1.2, 5, 10}
	// used to p95, p99
	Metrics.SetDuration([]float64{0.02, 0.1, 1.0, 5, 10})

	// set middleware for gin
	if config.Global.Server.Port == config.Global.Server.MetricsPort {
		go Metrics.Use(r)
		return
	}

	// use metric middleware without expose metric path
	Metrics.UseWithoutExposingEndpoint(r)
	rm := gin.Default()

	Metrics.Expose(rm)

	listen := fmt.Sprintf("%s:%d", config.Global.Server.Address, config.Global.Server.MetricsPort)
	logrus.Infof("Metrics server listen on: %s\n", listen)

	go rm.Run(listen)
}
