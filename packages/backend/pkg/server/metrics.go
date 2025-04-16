package server

import (
	"argocd-watcher/pkg/config"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/sirupsen/logrus"
)

func startMetrics(r *gin.Engine) {
	// get global Monitor object
	m := ginmetrics.GetMonitor()

	// +optional set metric path, default /debug/metrics
	m.SetMetricPath("/metrics")
	// +optional set slow time, default 5s
	m.SetSlowTime(10)
	// +optional set request duration, default {0.1, 0.3, 1.2, 5, 10}
	// used to p95, p99
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})

	// set middleware for gin
	if config.Global.Server.Port == config.Global.Server.MetricsPort {
		m.Use(r)
		return
	}

	// use metric middleware without expose metric path
	m.UseWithoutExposingEndpoint(r)
	listen := fmt.Sprintf("%s:%d", config.Global.Server.Address, config.Global.Server.MetricsPort)
	logrus.Infof("Metrics server listen on: %s\n", listen)

	rm := gin.Default()
	go rm.Run(listen)
	m.Expose(rm)
}
