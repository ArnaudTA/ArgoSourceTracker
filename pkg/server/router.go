package server

import (
	"argocd-watcher/pkg/config"
	"strconv"

	"github.com/gin-gonic/gin"
)

func favIcon(c *gin.Context) {
	c.File("favicon.ico")
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/favicon.ico", favIcon)
	r.GET("/apps", fetchApplications)
	r.GET("/origin/:instance", getApplicationOrigin)
	return r
}

func StartGin(cfg config.Config) {
	r := setupRouter()
	serverPort := strconv.Itoa(cfg.Server.Port)
	serverAddr := cfg.Server.Address
	r.Run(serverAddr + ":" + serverPort)
}
