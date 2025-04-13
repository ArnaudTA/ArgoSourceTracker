package parser

import (
	"argocd-watcher/pkg/registries"
	"strings"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/blang/semver/v4"
)

type ChartVersion struct {
	RepoURL  string   `json:"repoURL"`
	Status   string   `json:"status,omitempty"`
	Revision string   `json:"revision"`
	NewTags  []string `json:"newTags,omitempty"`
	Protocol 	string    `json:"type"`
}

func ParseSource(source v1alpha1.ApplicationSource, revision string) ChartVersion {
	chartVersion := ChartVersion{
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
		chartVersion.Status = "Obsolete"
	}
	return chartVersion
}

type ApplicationSummary struct {
	ChartVersions []ChartVersion `json:"chartVersions"`
	Status        string         `json:"Status"`
}

func ParseApplication(app v1alpha1.Application) []ChartVersion {
	chartVersions := []ChartVersion{}
	syncStatus := app.Status.Sync

	if len(syncStatus.Revisions) != 0 {
		for idx, source := range syncStatus.ComparedTo.Sources {
			if source.Chart == "" {
				continue
			}
			chartVersions = append(chartVersions, ParseSource(source, syncStatus.Revisions[idx]))
		}
	} else if syncStatus.Revision != "" && syncStatus.ComparedTo.Source.Size() == 0 {
		chartVersions = append(chartVersions, ParseSource(syncStatus.ComparedTo.Source, syncStatus.Revision))
	}

	return chartVersions
}
