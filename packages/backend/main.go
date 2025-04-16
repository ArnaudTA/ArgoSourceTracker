// @title ArgoSourceTracker API
// @version 1.0
// @description API simple pour lister les applications ArgoCD et suivre les versions des charts
// @BasePath /
package main

import (
	"argocd-watcher/pkg/application"
	"argocd-watcher/pkg/applicationset"
	"argocd-watcher/pkg/argocd"
	"argocd-watcher/pkg/config"
	"argocd-watcher/pkg/server"

	"github.com/sirupsen/logrus"
)

func main() {

	if err := config.LoadGlobal(); err != nil {
		logrus.Fatal(err)
	}
	logrus.SetLevel(logrus.WarnLevel)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	argocd.LoadArgoConf()

	appClient := argocd.GetClient()
	go application.Watch(appClient)
	go applicationset.Watch(appClient)
	server.StartGin()
}
