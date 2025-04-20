package application

import (
	"argocd-watcher/pkg/argocd/applicationset"
	"argocd-watcher/pkg/config"
	"argocd-watcher/pkg/registries"
	"argocd-watcher/pkg/types"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/blang/semver/v4"
)

var InstanceLabel string = "argocd.argoproj.io/instance"

func GenerateChartSummary(source *types.ApplicationSourceWithRevision) types.ChartSummary {
	version := semver.MustParse(source.Revision)
	status := "Unknown"
	index, err := registries.Search(source.Source.RepoURL)
	newTags := []string{}
	if err != nil {
		status = err.Error()
	} else {
		newTags = registries.GetGreaterTags(index, source.Source.RepoURL, source.Source.Chart, version)
		if len(newTags) > 0 {
			status = "Outdated"
		} else {
			status = "Up-to-date"
		}
	}
	return types.ChartSummary{
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
