package parser

import (
	"argocd-watcher/pkg/argocd/application"
	"argocd-watcher/pkg/registries"
	"strings"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/blang/semver/v4"
)

type Source struct {
	RepoURL  string   `json:"repoURL"`
	Status   string   `json:"status,omitempty"`
	Revision string   `json:"revision"`
	NewTags  []string `json:"newTags,omitempty"`
	Protocol string   `json:"type"`
	Chart    string   `json:"chart" binding:"required"`
}

func ParseSource(source v1alpha1.ApplicationSource, revision string) Source {
	sourceVersion := Source{
		RepoURL:  source.RepoURL,
		Revision: revision,
		Chart:    source.Chart,
	}
	semVerRevision := semver.MustParse(revision)
	if strings.HasPrefix(source.RepoURL, "https://") || strings.HasPrefix(source.RepoURL, "http://") {
		sourceVersion.NewTags = registries.GetGreaterTags(source.RepoURL, source.Chart, semVerRevision)
		sourceVersion.Protocol = "HTTP"
	} else {
		sourceVersion.Protocol = "OCI"
	}
	if len(sourceVersion.NewTags) == 0 {
		sourceVersion.Status = "Up-to-date"
	} else {
		sourceVersion.Status = "Outdated"
	}
	return sourceVersion
}

type ApplicationSummary struct {
	Charts   []Source `json:"charts" binding:"required"`
	Status   string   `json:"status,omitempty" binding:"required"`
	Instance string   `json:"instance,omitempty" binding:"required"`
}

func ParseApplication(app v1alpha1.Application) ApplicationSummary {
	charts := []Source{}
	syncStatus := app.Status.Sync

	if len(syncStatus.Revisions) != 0 {
		for idx, source := range syncStatus.ComparedTo.Sources {
			if source.Chart == "" {
				continue
			}
			charts = append(charts, ParseSource(source, syncStatus.Revisions[idx]))
		}
	} else if syncStatus.Revision != "" && syncStatus.ComparedTo.Source.Size() == 0 {
		charts = append(charts, ParseSource(syncStatus.ComparedTo.Source, syncStatus.Revision))
	}

	appSummary := ApplicationSummary{
		Charts: charts,
		Status: getApplicationStatus(charts),
	}

	if instance, ok := app.Labels[application.InstanceLabel]; ok {
		appSummary.Instance = instance
	}

	return appSummary
}

func getApplicationStatus(charts []Source) string {
	for _, chart := range charts {
		if chart.Status == "Outdated" {
			return "Outdated"
		}
	}
	return "Up-to-date"
}
