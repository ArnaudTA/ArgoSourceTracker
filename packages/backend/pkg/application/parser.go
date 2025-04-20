package application

import (
	// "argocd-watcher/pkg/applicationset"
	"argocd-watcher/pkg/applicationset"
	"argocd-watcher/pkg/config"
	"argocd-watcher/pkg/registries"
	"argocd-watcher/pkg/types"
	"fmt"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/blang/semver/v4"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var InstanceLabel string = "argocd.argoproj.io/instance"

func getApplicationStatus(charts []types.ChartSummary) types.ApplicationStatus {
	if len(charts) == 0 {
		return types.Ignored
	}
	for _, chart := range charts {
		if chart.Status == "Outdated" {
			return types.Outdated
		}
	}
	return types.UpToDate
}

func GenerateChartSummary(source types.ApplicationSourceWithRevision) types.ChartSummary {
	version := semver.MustParse(source.Revision)
	status := "Unknown"
	index, err := registries.Search(source.RepoURL)
	newTags := []string{}
	if err != nil {
		status = err.Error()
	} else {
		newTags = registries.GetGreaterTags(index, source.RepoURL, source.Chart, version)
		if len(newTags) > 0 {
			status = "Outdated"
		} else {
			status = "Up-to-date"
		}
	}
	return types.ChartSummary{
		RepoURL:  source.RepoURL,
		Status:   status,
		Revision: source.Revision,
		NewTags:  newTags,
		Chart:    source.Chart,
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
				app := owner.(Application)
				return &types.Parent{
					Kind:           "Application",
					Name:           value,
					Namespace:      metadata.Namespace,
					ApplicationUrl: getApplicationUrl(app.Resource.Namespace, app.Resource.Name),
				}, app.Resource.ObjectMeta.DeepCopy()
			}
		}
	}
	return nil, nil
}

func getApplicationUrl(namespace, name string) string {
	return fmt.Sprintf("%s/applications/%s/%s", config.Global.Argocd.Url, namespace, name)
}
