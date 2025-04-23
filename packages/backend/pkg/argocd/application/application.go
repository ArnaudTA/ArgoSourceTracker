package application

import (
	"argocd-watcher/pkg/config"
	"argocd-watcher/pkg/registries"
	"argocd-watcher/pkg/types"
	"fmt"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
)

type Application struct {
	Resource *v1alpha1.Application
	Sources  []*types.ApplicationSourceWithRevision
}

func (app *Application) GetSummary() types.Summary {
	charts := []types.ChartSummary{}
	for _, source := range app.Sources {
		if source.Source.Chart == "" {
			continue
		}
		charts = append(charts, GenerateChartSummary(source))
	}
	instance := app.Resource.Labels[InstanceLabel]

	return types.Summary{
		Charts:         charts,
		Status:         mostSevereStatus(charts),
		Instance:       instance,
		Name:           app.Resource.Name,
		Namespace:      app.Resource.Namespace,
		ApplicationUrl: app.getApplicationUrl(),
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

func (app *Application) ExtractSources() []*types.ApplicationSourceWithRevision {
	syncStatus := app.Resource.Status.Sync

	sources := []*types.ApplicationSourceWithRevision{}

	if len(syncStatus.Revisions) != 0 {
		for idx, source := range syncStatus.ComparedTo.Sources {
			if source.Chart == "" {
				continue
			}
			sources = append(sources, &types.ApplicationSourceWithRevision{
				Source:   source,
				Revision: syncStatus.Revisions[idx],
			})
		}
	} else if syncStatus.Revision != "" && syncStatus.ComparedTo.Source.Size() != 0 {
		sources = append(sources, &types.ApplicationSourceWithRevision{
			Source:   syncStatus.ComparedTo.Source,
			Revision: syncStatus.Revision,
		})
	}
	return sources
}

func (app *Application) Parse() {
	sources := app.ExtractSources()
	app.Sources = sources
	for _, source := range app.Sources {
		if source.Source.RepoURL == "" || source.Source.Chart == "" {
			continue
		}
		registries.Search(source.Source.RepoURL, source.Source.Chart)
	}
}

func (app *Application) getApplicationUrl() string {
	return fmt.Sprintf("%s/applications/%s/%s", config.Global.Argocd.Url, app.Resource.ObjectMeta.Namespace, app.Resource.ObjectMeta.Name)
}

func mostSevereStatus(charts []types.ChartSummary) types.ApplicationStatus {
	if len(charts) == 0 {
		return types.Ignored
	}
	worst := types.Ignored
	maxSeverity := types.Severity[worst]
	for _, chart := range charts {
		if sev, ok := types.Severity[chart.Status]; ok && sev > maxSeverity {
			worst = chart.Status
			maxSeverity = sev
		}
	}
	return worst
}
