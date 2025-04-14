package registries

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/blang/semver/v4"
	"gopkg.in/yaml.v3"
	"sync"
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

func StoreGet(registry string) (IndexFile, error) {
	cachedVersion, ok := Cache.Load(registry)
	if ok {
		fmt.Printf("Use cache for: %s\n", registry)
		return cachedVersion.(IndexFile), nil
	}

	fmt.Printf("Fetching: %s\n", registry)

	// Appel HTTP pour récupérer le fichier index.yaml
	resp, err := http.Get(registry + "/index.yaml")
	if err != nil {
		log.Fatalf("Erreur HTTP : %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Code HTTP inattendu : %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Erreur de lecture du body : %v", err)
		panic(err)
	}

	// Parsing YAML au format Helm
	index := IndexFile{}
	if err := yaml.Unmarshal(body, &index); err != nil {
		log.Fatalf("Erreur de parsing YAML : %v", err)
	}

	Cache.Store(registry, index)

	time.AfterFunc(300*time.Second, func() { StoreDeleteRegistry(registry) })

	return index, nil
}

func StoreDeleteRegistry(registry string) {
	Cache.Delete(registry)
}

func GetTags(registry, chart string) []string {
	tags := []string{}
	index, _ := StoreGet(registry)
	if entry, ok := index.Entries[chart]; ok {
		for _, version := range entry {
			tags = append(tags, version.Version)
		}
	}
	return tags
}

func GetGreaterTags(registry, chart string, minVersion semver.Version) []string {
	tags := []string{}
	index, err := StoreGet(registry)
	if err != nil {
		return tags
	}
	if entry, ok := index.Entries[chart]; ok {
		for _, version := range entry {
			candidateVersion, _ := semver.Parse(version.Version)
			if minVersion.LT(candidateVersion) {
				tags = append(tags, version.Version)
			}
		}
	}
	return tags
}
