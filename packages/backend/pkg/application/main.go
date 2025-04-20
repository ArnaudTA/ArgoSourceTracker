package application

import (
	"argocd-watcher/pkg/registries"
	"argocd-watcher/pkg/types"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
)

type Application struct {
	Resource *v1alpha1.Application
}

func (app *Application) GetSummary() types.Summary {
	charts := []types.ChartSummary{}
	sources := app.ExtractSources()
	for _, source := range sources {
		charts = append(charts, GenerateChartSummary(source))
	}
	instance := app.Resource.Labels[InstanceLabel]

	return types.Summary{
		Charts:         charts,
		Status:         getApplicationStatus(charts),
		Instance:       instance,
		Name:           app.Resource.Name,
		Namespace:      app.Resource.Namespace,
		ApplicationUrl: getApplicationUrl(app.Resource.Namespace, app.Resource.Name),
	}
}

func (app *Application) GetGenealogy() []*types.Parent {
	track := []*types.Parent{}
	metadata := &app.Resource.ObjectMeta

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

func (app *Application) ExtractSources() []types.ApplicationSourceWithRevision {
	syncStatus := app.Resource.Status.Sync

	sources := []types.ApplicationSourceWithRevision{}

	if len(syncStatus.Revisions) != 0 {
		for idx, source := range syncStatus.ComparedTo.Sources {
			if source.Chart == "" {
				continue
			}
			sources = append(sources, types.ApplicationSourceWithRevision{
				source,
				syncStatus.Revisions[idx],
			})
		}
	} else if syncStatus.Revision != "" && syncStatus.ComparedTo.Source.Size() == 0 {
		sources = append(sources, types.ApplicationSourceWithRevision{
			syncStatus.ComparedTo.Source,
			syncStatus.Revision,
		})
	}
	return sources
}

func (app *Application) Parse() {
	sources := app.ExtractSources()
	for _, source := range sources {
		if source.RepoURL == "" {
			continue
		}
		registries.Search(source.RepoURL)
		if source.Chart == "" {
			continue
		}
	}
}
