package application

import (
	"fmt"
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

var Store = map[string]v1alpha1.Application{}

func StoreGet(applicationName string) (v1alpha1.Application, error) {
	cachedApplication, ok := Store[applicationName]
	if ok {
		fmt.Printf("Use Application cache for: %s\n", applicationName)
		return cachedApplication, nil
	}

	fmt.Printf("\nFetching application: %s\n", applicationName)

	application, err := getApplication(applicationName)
	if err != nil {
		panic(err)
	}

	Store[applicationName] = *application

	time.AfterFunc(300*time.Second, func() { StoreDeleteApplication(applicationName) })

	return *application, nil
}

func StoreList() ([]v1alpha1.Application, error) {
	fmt.Printf("Fetching all applications\n")

	applications, err := getApplications()
	if err != nil {
		panic(err)
	}

	for _, app := range applications {
		Store[app.Name] = app
		time.AfterFunc(300*time.Second, func() { StoreDeleteApplication(app.Name) })
	}

	return applications, nil
}

func StoreDeleteApplication(applicationName string) {
	delete(Store, applicationName)
}
