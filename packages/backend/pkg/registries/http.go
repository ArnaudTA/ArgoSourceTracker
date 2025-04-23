package registries

import (
	"errors"
	"io"
	"net/http"

	"github.com/go-redis/cache/v9"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func getHttpChart(repositoryUrl, chartName string) (Chart, error) {
	var chart Chart
	entries, err := callHTTP(repositoryUrl)

	if err != nil && err != cache.ErrCacheMiss {
		return Chart{}, err
	}

	storeEntries(repositoryUrl, entries)

	key := generateKey(repositoryUrl, chartName)

	chart, err = ChartCache.Get(key)

	if err != nil && err != cache.ErrCacheMiss {
		return Chart{}, err
	}
	return chart, nil
}

func callHTTP(repo string) (Entries, error) {
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
