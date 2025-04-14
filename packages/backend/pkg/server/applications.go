package server

import (
	"argocd-watcher/pkg/argocd/application"
	"argocd-watcher/pkg/parser"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Liste les applications
// @Description Retourne la liste des applications et le rapport des versions
// @Tags Applications
// @Produce json
// @Success 200 {object} map[string]parser.ApplicationSummary
// @Param filter query string false "Filtre les applications"
// @Router /api/v1/apps [get]
func fetchApplications(c *gin.Context) {
	filter := c.DefaultQuery("filter", "standard")
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
		switch filter {
		case "standard":
			if len(appSummary.Charts) != 0 {
				result[app.Name] = appSummary
			}
		case "outdated":
			if appSummary.Status == "Outdated" {
				result[app.Name] = appSummary
			}
		case "all":
			result[app.Name] = appSummary
		}
	}

	c.JSON(http.StatusOK, result)
}

// @Summary Remonte l'origine d'une application
// @Description Liste les applications et applications qui ménent à cette application
// @Tags Track Origin
// @Produce json
// @Success 200 {array} application.TrackRecord
// @Failure 400 {object} error
// @Param application path string true "Application cible"
// @Router /api/v1/apps/{application}/origin [get]
func getApplicationOrigin(c *gin.Context) {
	name := c.Param("application")
	if name == "" {
		c.Abort()
	}
	appTrack := application.GetApplicationTrack(name)

	c.JSON(http.StatusOK, appTrack)
}
