package registries

import (
	"argocd-watcher/pkg/cache"
	"bytes"
	"encoding/gob"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
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
type IndexFile struct {
	Entries map[string]ChartVersions `json:"entries"`
}

func Search(registry string) (*IndexFile, error) {
	if !strings.HasPrefix(registry, "https://") && !strings.HasPrefix(registry, "http://") {
		return nil, errors.New("invalid url")
	}
	redisKey := generateKey(registry)

	resultBytes, err := cache.Load(redisKey)
	if err == redis.Nil {
		logrus.Debugf("Fetching: %s\n", registry)

		// Appel HTTP pour récupérer le fichier index.yaml
		resp, err := http.Get(registry + "/index.yaml")
		if err != nil {
			logrus.Fatalf("Erreur HTTP : %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			logrus.Fatalf("Code HTTP inattendu : %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logrus.Fatalf("Erreur de lecture du body : %v", err)
			panic(err)
		}

		// Parsing YAML au format Helm
		index := IndexFile{}
		if err := yaml.Unmarshal(body, &index); err != nil {
			logrus.Fatalf("Erreur de parsing YAML : %v", err)
		}

		var buffer bytes.Buffer
		enc := gob.NewEncoder(&buffer)
		err = enc.Encode(index)
		if err != nil {
			logrus.Fatalf("Erreur d'encodage gob: %v", err)
		}
		indexGOB := buffer.Bytes()
		cache.Store(redisKey, indexGOB, 0)

		time.AfterFunc(300*time.Second, func() { Delete(registry) })

		return &index, nil
	}
	if err != nil {
		logrus.WithError(err)
		return nil, err
	}
	logrus.Debugf("Use cache for: %s\n", registry)
	var cachedIndex *IndexFile
	dec := gob.NewDecoder(bytes.NewReader(resultBytes))
	err = dec.Decode(&cachedIndex)
	if err != nil {
		logrus.Errorf("Erreur de décodage gob: %v", err)
	}
	return cachedIndex, nil
}

func Delete(registry string) {
	cache.Delete(generateKey(registry))
}

func generateKey(registry string) string {
	registry = strings.TrimSuffix(registry, "/")
	return "ast/registry/" + registry
}
