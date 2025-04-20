package server

import (
	"argocd-watcher/pkg/application"
	"argocd-watcher/pkg/types"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Summary Liste les applications
// @Description Retourne la liste des applications et le rapport des versions
// @Tags Applications
// @Produce json
// @Success 200 {array} application.ApplicationSummary
// @Param name query string false "Name to search"
// @Param filter query string false "Filtre les applications"
// @Router /api/v1/apps [get]
func fetchApplications(c *gin.Context) {
	filterQuery := c.DefaultQuery("filter", "")
	nameQuery := c.DefaultQuery("name", "")

	summaries := []types.Summary{}
	application.AppCache.Range(func(key, value any) bool {
		name := fmt.Sprint(key)
		if nameQuery != "" && !strings.Contains(name, nameQuery) {
			return true
		}
		app := value.(application.Application)
		summary := app.GetSummary()
		switch filterQuery {
		case "all":
			summaries = append(summaries, summary)
		case "outdated":
			if summary.Status == types.Outdated {
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

	result, ok := application.AppCache.Load(key)
	if !ok {
		c.AbortWithStatus(404)
	}
	app := result.(application.Application)
	summary := app.GetSummary()

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
	app := item.(application.Application)
	appTrack := app.GetGenealogy()

	c.JSON(http.StatusOK, appTrack)
}
