package server

import (
	"github.com/gin-gonic/gin"
)

func favIcon (c *gin.Context) {
	c.File("favicon.ico")
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/favicon.ico", favIcon)
	r.GET("/apps", fetchApplications)
	r.GET("/cache", getCache)
	r.GET("/cache/keys", getCacheKeys)
	r.GET("/invalidate", invalidateCache)
	return r
}

func StartGin() {
	r := setupRouter()
	r.Run(":8080")
}