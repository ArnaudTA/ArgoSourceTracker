package parser

import (
	"argocd-watcher/pkg/argocd/application"
	"argocd-watcher/pkg/registries"
	"strings"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/blang/semver/v4"
)

type Chart struct {
	RepoURL  string   `json:"repoURL"`
	Status   string   `json:"status,omitempty"`
	Revision string   `json:"revision"`
	NewTags  []string `json:"newTags,omitempty"`
	Protocol string   `json:"type"`
}

func ParseSource(source v1alpha1.ApplicationSource, revision string) Chart {
	chartVersion := Chart{
		RepoURL:  source.RepoURL,
		Revision: revision,
	}
	semVerRevision := semver.MustParse(revision)
	if strings.HasPrefix(source.RepoURL, "https://") || strings.HasPrefix(source.RepoURL, "http://") {
		chartVersion.NewTags = registries.GetGreaterTags(source.RepoURL, source.Chart, semVerRevision)
		chartVersion.Protocol = "HTTP"
	} else {
		chartVersion.Protocol = "OCI"
	}
	if len(chartVersion.NewTags) == 0 {
		chartVersion.Status = "Up-to-date"
	} else {
		chartVersion.Status = "Outdated"
	}
	return chartVersion
}

type ApplicationSummary struct {
	Charts   []Chart `json:"charts"`
	Status   string  `json:"status,omitempty"`
	Instance string  `json:"instance,omitempty"`
}

func ParseApplication(app v1alpha1.Application) ApplicationSummary {
	charts := []Chart{}
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

func getApplicationStatus(charts []Chart) string {
	for _, chart := range charts {
		if chart.Status == "Outdated" {
			return "Outdated"
		}
	}
	return "Up-to-date"
}
