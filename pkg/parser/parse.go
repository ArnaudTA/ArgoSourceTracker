package parser

import (
	"argocd-watcher/pkg/chart"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/blang/semver/v4"
)

type ChartVersion struct {
	Application string   `json:"application,omitempty"`
	RepoURL     string   `json:"repoURL"`
	Status      string   `json:"status,omitempty"`
	Revision    string   `json:"revision"`
	NewTags     []string `json:"newTags,omitempty"`
}

func ParseSource(source v1alpha1.ApplicationSource, revision string, appName string) ChartVersion {
	chartVersion := ChartVersion{
		Application: appName,
		RepoURL:     source.RepoURL,
		Status:      "Unknown",
		Revision:    revision,
	}
	semVerRevision := semver.MustParse(revision)
	chartVersion.NewTags = chart.GetGreaterTags(source.RepoURL, source.Chart, semVerRevision)
	return chartVersion
}

func ParseApplication(app v1alpha1.Application) []ChartVersion {
	appName := app.Name
	appSources := []ChartVersion{}
	syncStatus := app.Status.Sync
	if syncStatus.Size() == 0 {
		return appSources
	}
	if len(syncStatus.Revisions) != 0 {
		for idx, source := range syncStatus.ComparedTo.Sources {
			if source.Chart == "" {
				continue
			}
			appSources = append(appSources, ParseSource(source, syncStatus.Revisions[idx], appName))
		}
	} else if syncStatus.Revision != "" && syncStatus.ComparedTo.Source.Size() == 0 {
		appSources = append(appSources, ParseSource(syncStatus.ComparedTo.Source, syncStatus.Revision, appName))
	}
	return appSources
}
