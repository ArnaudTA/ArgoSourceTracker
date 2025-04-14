package server

import (
	"argocd-watcher/pkg/config"
	"net/http"
	"strconv"

	_ "argocd-watcher/docs" // 👈 important pour init Swagger

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Check struct {
	Status string `json:"status,omitempty" binding:"required"`
}

// @Summary Status
// @Description Retourne le status de l'application
// @Tags Healthcheck
// @Produce json
// @Success 200 {object} Check
// @Router /api/v1/health [get]
func health(c *gin.Context) {
	c.JSON(200, Check{ Status: "OK" })
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	// Enable CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// API routes
	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", health)
		v1.GET("/apps", fetchApplications)
		v1.GET("/apps/:application", fetchApplication)
		v1.GET("/apps/:application/origin", getApplicationOrigin)
	}

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Serve static files from the Vue app
	r.StaticFS("/ui", http.Dir("static"))

	// Serve static files from the Vue app
	r.StaticFS("/assets", http.Dir("static/assets"))

	// Serve static files from the Vue app
	r.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/ui")
	})

	r.Run(":8080")

	return r
}

func StartGin() {
	r := setupRouter()
	serverPort := strconv.Itoa(config.Global.Server.Port)
	serverAddr := config.Global.Server.Address
	r.Run(serverAddr + ":" + serverPort)
}
