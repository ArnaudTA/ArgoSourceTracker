package registries

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-redis/cache/v9"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
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

func Search(registryUrl, chartName string) (Chart, error) {
	if (!strings.HasPrefix(registryUrl, "https://") && !strings.HasPrefix(registryUrl, "http://")) || strings.HasSuffix(registryUrl, ".git") {
		return Chart{}, errors.ErrUnsupported
	}

	key := generateKey(registryUrl, chartName)
	cached, err := ChartCache.Get(key)

	if err == nil {
		return cached, nil
	}
	entries, err := CallHTTP(registryUrl)

	if err != nil && err != cache.ErrCacheMiss {
		return Chart{}, err
	}

	storeEntries(registryUrl, entries)
	cached, err = ChartCache.Get(key)

	if err != nil && err != cache.ErrCacheMiss {
		return Chart{}, err
	}
	return cached, nil
	// if (!strings.HasPrefix(registryUrl, "https://") && !strings.HasPrefix(registryUrl, "http://")) || strings.HasSuffix(registryUrl, ".git") {
	// 	return nil, errors.New("Don't support OCI registry yet")
	// }

	// if err != nil {
	// 	logrus.WithError(err)
	// 	return nil, err
	// }
	// logrus.Debugf("Use cache for: %s\n", registryUrl)
}

func generateKey(repository, chart string) string {
	repository = strings.TrimSuffix(repository, "/")
	return fmt.Sprintf("ast/registry/%s/%s", repository, chart)
}

func CallHTTP(repo string) (Entries, error) {
	logrus.Debugf("Fetching: %s\n", repo)

	// Appel HTTP pour récupérer le fichier index.yaml
	resp, err := http.Get(repo + "/index.yaml")
	if err != nil {
		logrus.Fatalf("Erreur HTTP : %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("Code HTTP inattendu : %d, for %s", resp.StatusCode, repo)
		return nil, errors.New("unexpected http code")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Erreur de lecture du body : %v", err)
		return nil, err
	}

	// Parsing YAML au format Helm
	index := IndexFile{}
	if err := yaml.Unmarshal(body, &index); err != nil {
		logrus.Errorf("Erreur de parsing YAML : %v", err)
		return nil, err
	}

	return index.Entries, nil
}

func storeEntries(repo string, entries Entries) {
	for chartName, entry := range entries {
		chart := Chart{}
		for _, metadata := range entry {
			chart.Tags = append(chart.Tags, metadata.Version)
		}
		key := generateKey(repo, chartName)
		ChartCache.Set(key, chart)
	}
}
