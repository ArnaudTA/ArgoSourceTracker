package applicationset

import (
	"fmt"
	"sync"
	"time"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
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
type IndexFile struct {
	Entries map[string]ChartVersions `json:"entries"`
}

var Cache = sync.Map{}

func StoreGet(applicationsetName string) (v1alpha1.ApplicationSet, error) {
	cachedApplicationset, ok := Cache.Load(applicationsetName)
	if ok {
		fmt.Printf("Use Application cache for: %s\n", applicationsetName)
		return cachedApplicationset.(v1alpha1.ApplicationSet), nil
	}

	fmt.Printf("\nFetching application: %s\n", applicationsetName)

	applicationset, err := getApplicationset(applicationsetName)
	if err != nil {
		panic(err)
	}

	Cache.Store(applicationsetName, *applicationset)

	time.AfterFunc(300*time.Second, func() { storeDeleteApplicationset(applicationsetName) })

	return *applicationset, nil
}

func StoreList() ([]v1alpha1.ApplicationSet, error) {
	fmt.Printf("Fetching all applications\n")

	applicationsets, err := getApplicationsets()
	if err != nil {
		panic(err)
	}

	for _, appset := range applicationsets {
		Cache.Store(appset.Name, appset)
		time.AfterFunc(300*time.Second, func() { storeDeleteApplicationset(appset.Name) })
	}

	return applicationsets, nil
}

func storeDeleteApplicationset(applicationName string) {
	Cache.Delete(applicationName)
}
