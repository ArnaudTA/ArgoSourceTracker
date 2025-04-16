package server

import (
	"argocd-watcher/pkg/application"
	"net/http"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/gin-gonic/gin"
)

// @Summary Liste les applications
// @Description Retourne la liste des applications et le rapport des versions
// @Tags Applications
// @Produce json
// @Success 200 {array} application.ApplicationSummary
// @Param filter query string false "Filtre les applications"
// @Router /api/v1/apps [get]
func fetchApplications(c *gin.Context) {
	filter := c.DefaultQuery("filter", "")

	summaries := []application.ApplicationSummary{}
	application.AppCache.Range(func(key, value any) bool {
		summary := application.GenerateApplicationSummary(value.(*v1alpha1.Application))
		switch filter {
		case "all":
			summaries = append(summaries, summary)
		case "outdated":
			if summary.Status == "Outdated" {
				summaries = append(summaries, summary)
			}
		default:
			if len(summary.Charts) > 0 {
				summaries = append(summaries, summary)
			}
		}

		return true
	})

	c.JSON(http.StatusOK, summaries)
}

// @Summary Récupe une application
// @Description Retourne application et le rapport de versions
// @Tags Applications
// @Produce json
// @Success 200 {object} application.ApplicationSummary
// @Param name path string true "Application cible"
// @Param namespace path string true "Namespace cible"
// @Router /api/v1/apps/{namespace}/{name} [get]
func fetchApplication(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	key := namespace + "/" + name

	app, ok := application.AppCache.Load(key)
	if !ok {
		c.AbortWithStatus(404)
	}
	summary := application.GenerateApplicationSummary(app.(*v1alpha1.Application))

	c.JSON(http.StatusOK, summary)
}

// @Summary Remonte l'origine d'une application
// @Description Liste les applications et applications qui ménent à cette application
// @Tags Track Origin
// @Produce json
// @Success 200 {array} application.GenealogyItem
// @Failure 400 {object} error
// @Param name path string true "Application cible"
// @Param namespace path string true "Namespace cible"
// @Router /api/v1/apps/{namespace}/{name}/origin [get]
func getApplicationOrigin(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	key := namespace + "/" + name

	if name == "" {
		c.Abort()
	}
	item, ok := application.AppCache.Load(key)
	if !ok {
		c.AbortWithStatus(404)
		return
	}
	appTrack := application.GetApplicationTrack(item.(*v1alpha1.Application))

	c.JSON(http.StatusOK, appTrack)
}
