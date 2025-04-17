package application

import (
	// "argocd-watcher/pkg/applicationset"
	"argocd-watcher/pkg/applicationset"
	"argocd-watcher/pkg/config"
	"argocd-watcher/pkg/registries"
	"fmt"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/blang/semver/v4"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var InstanceLabel string = "argocd.argoproj.io/instance"

func ParseApplication(app *v1alpha1.Application) {
	sources := ExtractSources(app)
	for _, source := range sources {
		registries.StoreGet(source.RepoURL)
	}
}

func getApplicationStatus(charts []ChartSummary) ApplicationStatus {
	if len(charts) == 0 {
		return ApplicationStatus(Ignored)
	}
	for _, chart := range charts {
		if chart.Status == "Outdated" {
			return ApplicationStatus(Outdated)
		}
	}
	return ApplicationStatus(UpToDate)
}

func GenerateApplicationSummary(app *v1alpha1.Application) ApplicationSummary {
	charts := []ChartSummary{}
	sources := ExtractSources(app)
	for _, source := range sources {
		charts = append(charts, GenerateChartSummary(source))
	}
	instance := app.Labels[InstanceLabel]

	return ApplicationSummary{
		Charts:    charts,
		Status:    getApplicationStatus(charts),
		Instance:  instance,
		Name:      app.Name,
		Namespace: app.Namespace,
		ApplicationUrl: getApplicationUrl(app.Namespace, app.Name),
	}
}

func GenerateChartSummary(source ApplicationSourceWithRevision) ChartSummary {
	version := semver.MustParse(source.revision)
	status := "Unknown"
	index, err := registries.StoreGet(source.RepoURL)
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
	return ChartSummary{
		RepoURL:  source.RepoURL,
		Status:   status,
		Revision: source.revision,
		NewTags:  newTags,
		Chart:    source.Chart,
	}
}

func ExtractSources(app *v1alpha1.Application) []ApplicationSourceWithRevision {
	syncStatus := app.Status.Sync

	sources := []ApplicationSourceWithRevision{}

	if len(syncStatus.Revisions) != 0 {
		for idx, source := range syncStatus.ComparedTo.Sources {
			if source.Chart == "" {
				continue
			}
			sources = append(sources, ApplicationSourceWithRevision{
				source,
				syncStatus.Revisions[idx],
			})
		}
	} else if syncStatus.Revision != "" && syncStatus.ComparedTo.Source.Size() == 0 {
		sources = append(sources, ApplicationSourceWithRevision{
			syncStatus.ComparedTo.Source,
			syncStatus.Revision,
		})
	}
	return sources
}

func GetApplicationTrack(application *v1alpha1.Application) []*GenealogyItem {
	track := []*GenealogyItem{}
	metadata := &application.ObjectMeta

	for i := 0; i < 10; i++ {
		item, parentMetadata := findParent(metadata)
		if item == nil {
			break
		}
		if parentMetadata == nil {
			item.ErrorMessage = "Can't find parent resource"
			break
		}
		track = append(track, item)
		metadata = parentMetadata
	}
	return track
}

type PreviousResource struct {
	Kind  string
	Name  string
	Found bool
}

func findParent(metadata *metav1.ObjectMeta) (*GenealogyItem, *metav1.ObjectMeta) {
	if metadata.OwnerReferences != nil {
		for _, ref := range metadata.DeepCopy().OwnerReferences {
			if ref.Kind == "ApplicationSet" {
				key := config.Global.Argocd.Namespace + "/" + ref.Name
				owner, ok := applicationset.AppSetCache.Load(key)
				if !ok {
					return &GenealogyItem{
						Kind:      ref.Kind,
						Name:      ref.Name,
						Namespace: metadata.Namespace,
					}, nil
				}
				appset := owner.(*v1alpha1.ApplicationSet)
				return &GenealogyItem{
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
					return &GenealogyItem{
						Kind:      "Application",
						Name:      value,
						Namespace: metadata.Namespace,
					}, nil
				}
				app := owner.(*v1alpha1.Application)
				return &GenealogyItem{
					Kind:           "Application",
					Name:           value,
					Namespace:      metadata.Namespace,
					ApplicationUrl: getApplicationUrl(app.Namespace, app.Name),
				}, app.ObjectMeta.DeepCopy()
			}
		}
	}
	return nil, nil
}


func getApplicationUrl(namespace, name string) string {
	return fmt.Sprintf("%s/applications/%s/%s", config.Global.Argocd.Url, namespace, name)
}