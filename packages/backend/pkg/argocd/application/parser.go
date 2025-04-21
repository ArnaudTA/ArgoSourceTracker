package application

import (
	"argocd-watcher/pkg/argocd/applicationset"
	"argocd-watcher/pkg/config"
	"argocd-watcher/pkg/registries"
	"argocd-watcher/pkg/types"
	"fmt"
	"sort"
	"strings"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/blang/semver/v4"
)

var InstanceLabel string = "argocd.argoproj.io/instance"

func GenerateChartSummary(source *types.ApplicationSourceWithRevision) types.ChartSummary {
	var status types.ApplicationStatus
	errorMessage := ""
	version, err := semver.Parse(source.Revision)
	newTags := []string{}
	if err != nil {
		status = types.Ignored
		errorMessage = "Invalid semver"
	} else {
		index, err := registries.Search(source.Source.RepoURL)
		if err != nil {
			status = types.Error
			errorMessage = err.Error()
		} else {
			newTags = registries.GetGreaterTags(index, source.Source.RepoURL, source.Source.Chart, version)
			if len(newTags) > 0 {
				status = "Outdated"
			} else {
				status = "Up-to-date"
			}
		}
	}
	return types.ChartSummary{
		Error:    errorMessage,
		RepoURL:  source.Source.RepoURL,
		Status:   status,
		Revision: source.Revision,
		NewTags:  newTags,
		Chart:    source.Source.Chart,
	}
}

type PreviousResource struct {
	Kind  string
	Name  string
	Found bool
}

func findParent(metadata *metav1.ObjectMeta) (*types.Parent, *metav1.ObjectMeta) {
	if metadata.OwnerReferences != nil {
		for _, ref := range metadata.DeepCopy().OwnerReferences {
			if ref.Kind == "ApplicationSet" {
				key := config.Global.Argocd.Namespace + "/" + ref.Name
				owner, ok := applicationset.AppSetCache.Load(key)
				if !ok {
					return &types.Parent{
						Kind:      ref.Kind,
						Name:      ref.Name,
						Namespace: metadata.Namespace,
					}, nil
				}
				appset := owner.(*v1alpha1.ApplicationSet)
				return &types.Parent{
					Kind:      ref.Kind,
					Name:      ref.Name,
					Namespace: metadata.Namespace,
				}, appset.ObjectMeta.DeepCopy()
			}
		}
	}
	if metadata.Labels != nil {
		for key, value := range metadata.Labels {
			if key == config.Global.Argocd.InstanceLabelKey {
				key := config.Global.Argocd.Namespace + "/" + value
				owner, ok := AppCache.Load(key)
				if !ok {
					return &types.Parent{
						Kind:      "Application",
						Name:      value,
						Namespace: metadata.Namespace,
					}, nil
				}
				app := owner.(*Application)
				return &types.Parent{
					Kind:           "Application",
					Name:           value,
					Namespace:      metadata.Namespace,
					ApplicationUrl: app.getApplicationUrl(),
				}, app.Resource.ObjectMeta.DeepCopy()
			}
		}
	}
	return nil, nil
}

func List(statusFilter []string, name string, offset, limit int) types.ListApplicationRes {

	resp := types.ListApplicationRes{
		Items: []types.Summary{},
		Stats: types.AppStats{
			types.Error: 0,
			types.Ignored: 0,
			types.Outdated: 0,
			types.UpToDate: 0,
		},
	}
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

	keys := getKeys(&AppCache)
	sort.Strings(keys)
	offSetted := 0
	for _, key := range keys {
		appCached, ok := AppCache.Load(key)
		if !ok {
			continue
		}
		app := appCached.(*Application)
		summary := app.GetSummary()
		resp.Stats[summary.Status]++
		if len(statusFilter) > 0 && statusFilter[0] != "" && !isEnumInList(summary.Status, statusFilter) {
			continue
		}
		if name != "" && !strings.Contains(key, name) {
			continue
		}
		if offSetted < offset {
			offSetted++
			total++
			continue
		}
		if len(resp.Items) >= limit {
			total++
			continue
		}
		resp.Items = append(resp.Items, summary)
		offSetted++
		total++

	}

	AppCache.Range(func(key, value any) bool {
		appName := fmt.Sprint(key)
		if name != "" && !strings.Contains(appName, name) {
			return true
		}

		return true
	})
	return resp
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