// @title ArgoSourceTracker API
// @version 1.0
// @description API simple pour lister les applications ArgoCD et suivre les versions des charts
// @BasePath /
package main

import (
	"argocd-watcher/pkg/argocd"
	"argocd-watcher/pkg/argocd/application"
	"argocd-watcher/pkg/argocd/applicationset"
	"argocd-watcher/pkg/cache"
	"argocd-watcher/pkg/config"
	"argocd-watcher/pkg/metrics"
	"argocd-watcher/pkg/server"

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
	cache.GetClient()
	metrics.Init()
	server.StartGin()
}
