package registries

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Metadata for a Chart file. This models the structure of a Chart.yaml file.
//
// Spec: https://k8s.io/helm/blob/master/docs/design/chart_format.md#the-chart-file
type Metadata struct {
	// A SemVer 2 conformant version string of the chart
	Version string `json:"version,omitempty"`
	// The name of the chart
	Name string `json:"name,omitempty"`
}

// ChartVersions is a list of versioned chart references.
// Implements a sorter on Version.
type ChartVersions []Metadata

type Entries map[string]ChartVersions
type IndexFile struct {
	Entries `json:"entries"`
}

func Search(repositoryUrl, chartName string) (Chart, error) {
	key := generateKey(repositoryUrl, chartName)
	chart, err := ChartCache.Get(key)

	if err == nil {
		return chart, nil
	}

	if strings.HasSuffix(repositoryUrl, ".git") {
		return Chart{}, errors.ErrUnsupported
	}

	if ok, _ := regexp.MatchString("http(s)?://.*", repositoryUrl); ok {
		return getHttpChart(repositoryUrl, chartName)
	}

	// Means it's an OCI registry
	return getOciChart(repositoryUrl, chartName)
}

func generateKey(repository, chart string) string {
	repository = strings.TrimSuffix(repository, "/")
	return fmt.Sprintf("ast/registry/%s/%s", repository, chart)
}
