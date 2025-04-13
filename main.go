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
