package server

import (
	"argocd-watcher/pkg/argocd/application"
	"argocd-watcher/pkg/types"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

// @Summary Liste les applications
// @Description Retourne la liste des applications et le rapport des versions
// @Tags Applications
// @Produce json
// @Success 200 {array} types.Summary
// @Param name query string false "Name to search"
// @Param offset query int64 false "Number of elements to skip, default: 0"
// @Param limit query int64 false "Number of elements to return, default: 10"
// @Param filter query string false "Filtre les applications"
// @Router /api/v1/apps [get]
func listApplications(c *gin.Context) {
	statusQuery := c.DefaultQuery("filter", "")
	offsetQuery := c.DefaultQuery("offset", "")
	limitQuery := c.DefaultQuery("limit", "")
	nameQuery := c.DefaultQuery("name", "")

	statusFilter := strings.Split(statusQuery, ",")
	total := 0
	if len(statusFilter) == 1 && statusFilter[0] == "" {
		statusList := []types.ApplicationStatus{
			types.Outdated,
		}
		statusFilter := []string{}
		for _, status := range statusList {
			statusFilter = append(statusFilter, string(status))
		}
	}

	offset := stringToInt(offsetQuery, 0)
	limit := stringToInt(limitQuery, 10)
	summaries := []types.Summary{}

	keys := getKeys(&application.AppCache)
	sort.Strings(keys)
	offSetted := 0
	for _, key := range keys {
		appCached, ok := application.AppCache.Load(key)
		if !ok {
			continue
		}
		app := appCached.(*application.Application)
		summary := app.GetSummary()
		if len(statusFilter) > 0 && statusFilter[0] != "" && !isEnumInList(summary.Status, statusFilter) {
			continue
		}
		if nameQuery != "" && !strings.Contains(key, nameQuery) {
			continue
		}
		if offSetted < offset {
			offSetted++
			total++
			continue
		}
		if len(summaries) >= limit {
			total++
			continue
		}
		summaries = append(summaries, summary)
		offSetted++
		total++

	}

	application.AppCache.Range(func(key, value any) bool {
		name := fmt.Sprint(key)
		if nameQuery != "" && !strings.Contains(name, nameQuery) {
			return true
		}

		return true
	})
	c.Header("x-total", fmt.Sprintf("%d", total))
	c.Header("x-offset", fmt.Sprintf("%d", offset))
	c.JSON(http.StatusOK, summaries)
}

// @Summary Récupe une application
// @Description Retourne application et le rapport de versions
// @Tags Applications
// @Produce json
// @Success 200 {object} types.Summary
// @Header       200  {string}  x-total     "Total of items available"
// @Header       200  {string}  x-offset     "Return the offset you provided"
// @Param name path string true "Application cible"
// @Param namespace path string true "Namespace cible"
// @Router /api/v1/apps/{namespace}/{name} [get]
func getApplication(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	key := namespace + "/" + name

	result, ok := application.AppCache.Load(key)
	if !ok {
		c.AbortWithStatus(404)
	}
	app := result.(*application.Application)
	summary := app.GetSummary()

	c.JSON(http.StatusOK, summary)
}

// @Summary Remonte l'origine d'une application
// @Description Liste les applications et applications qui ménent à cette application
// @Tags Track Origin
// @Produce json
// @Success 200 {array} types.Parent
// @Failure 400 {object} error
// @Param name path string true "Application cible"
// @Param namespace path string true "Namespace cible"
// @Router /api/v1/apps/{namespace}/{name}/origin [get]
func getApplicationGenealogy(c *gin.Context) {
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
	app := item.(*application.Application)
	appTrack := app.GetGenealogy()

	c.JSON(http.StatusOK, appTrack)
}

func getKeys(m *sync.Map) []string {
	keys := []string{}

	m.Range(func(key, value any) bool {
		strKey := key.(string)
		keys = append(keys, strKey)
		return true // continuer à itérer
	})

	return keys
}

// Fonction pour vérifier si une valeur de l'enum est dans une liste de strings
func isEnumInList(val types.ApplicationStatus, list []string) bool {
	for _, item := range list {
		if string(val) == item {
			return true
		}
	}
	return false
}

func stringToInt(str string, defaultValue int) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return defaultValue
	}
	return i
}
