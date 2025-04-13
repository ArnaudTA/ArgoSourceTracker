package server

import (
	"fmt"

	"argocd-watcher/pkg/registries"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

func getCache(c *gin.Context) {
	registry := c.DefaultQuery("registry", "")
	fmt.Println(registry)
	if index, err := registries.StoreGet(registry); err != nil {
		body, _ := yaml.Marshal(index)
		place := []byte{72}
		c.Data(200, "string", append(body, place...))
	}
}

func getCacheKeys(c *gin.Context) {
	registry := c.DefaultQuery("registry", "")
	fmt.Println(registry)
	if index, err := registries.StoreGet(registry); err != nil {
		body, _ := yaml.Marshal(index)

		c.Data(200, "string", body)
	}
}

func invalidateCache(c *gin.Context) {
	registry := c.DefaultQuery("registry", "")
	fmt.Println(registry)
	registries.StoreInvalidate(registry)
	var result []gin.H
	c.JSON(200, result)
}
