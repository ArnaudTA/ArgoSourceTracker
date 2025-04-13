package server

import (
	"argocd-watcher/pkg/argocd"
	"argocd-watcher/pkg/parser"
	"net/http"

	"github.com/gin-gonic/gin"
)

func fetchApplications(c *gin.Context) {

	// Liste des applications ArgoCD
	applications, err := argocd.GetApplications()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Créer une liste simplifiée des applications
	result := map[string]parser.ApplicationSummary{}
	appList := applications.Items
	for _, app := range appList {
		appSummary := parser.ParseApplication(app)
		if len(appSummary.Charts) != 0 {
			result[app.Name] = appSummary
		}
	}

	c.JSON(http.StatusOK, result)
}
