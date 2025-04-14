// @title ArgoSourceTracker API
// @version 1.0
// @description API simple pour lister les applications ArgoCD et suivre les versions des charts
// @host localhost:8080
// @BasePath /
package main

import (
	"argocd-watcher/pkg/argocd"
	"argocd-watcher/pkg/config"
	"argocd-watcher/pkg/server"
	"log"
)

func main() {
	if err := config.LoadGlobal(); err != nil {
		log.Fatal(err)
	}
	argocd.InitClient()
	argocd.LoadArgoConf()
	server.StartGin()
}
