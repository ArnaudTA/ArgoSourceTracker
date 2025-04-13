package server

import (
	"argocd-watcher/pkg/argocd/application"
	"argocd-watcher/pkg/parser"
	"net/http"

	"github.com/gin-gonic/gin"
)

func fetchApplications(c *gin.Context) {
	// Liste des applications ArgoCD
	applications, err := application.StoreList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Créer une liste simplifiée des applications
	result := map[string]parser.ApplicationSummary{}
	appList := applications
	for _, app := range appList {
		appSummary := parser.ParseApplication(app)
		if len(appSummary.Charts) != 0 {
			result[app.Name] = appSummary
		}
	}

	c.JSON(http.StatusOK, result)
}

func getApplicationOrigin(c *gin.Context) {
	instance := c.Param("instance")

	appTrack := application.GetApplicationTrack(instance)

	c.JSON(http.StatusOK, appTrack)
}
