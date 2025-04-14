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
	"fmt"
	"log"
)

func main() {
	var cfg config.Config
	err := config.Load(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(cfg)
	argocd.InitClient(cfg)
	server.StartGin(cfg)
}
