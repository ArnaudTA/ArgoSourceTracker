// @title ChartSentinel API
// @version 1.0
// @description Track ArgoCD applications and charts version
// @BasePath /
package main

import (
	"github.com/cableship/chart-sentinel/pkg/argocd"
	"github.com/cableship/chart-sentinel/pkg/argocd/application"
	"github.com/cableship/chart-sentinel/pkg/argocd/applicationset"
	"github.com/cableship/chart-sentinel/pkg/config"
	"github.com/cableship/chart-sentinel/pkg/metrics"
	"github.com/cableship/chart-sentinel/pkg/registries"
	"github.com/cableship/chart-sentinel/pkg/server"

	"github.com/sirupsen/logrus"
)

func main() {
	if err := config.LoadGlobal(); err != nil {
		logrus.Fatal(err)
	}
	logrus.SetLevel(logrus.Level(config.Global.Log.Level))
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	argocd.LoadArgoConf()

	appClient := argocd.GetClient()
	go application.Watch(appClient)
	go applicationset.Watch(appClient)
	registries.Init()
	metrics.Init()
	server.StartGin()
}
