package registries

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"oras.land/oras-go/v2/registry/remote"
)

func getOciChart(repositoryUrl, chartName string) (Chart, error) {
	var chart Chart
	tags, err := callOCI(repositoryUrl, chartName)
	if err != nil {
		return chart, err
	}
	logrus.Warnln(tags)
	chart.Tags = tags
	key := generateKey(repositoryUrl, chartName)
	ChartCache.Set(key, chart)
	return chart, nil
}

func callOCI(repositoryUrl, chartName string) ([]string, error) {
	repositoryUrl = strings.TrimSuffix(repositoryUrl, "/")
	repositoryUrl = strings.TrimPrefix(repositoryUrl, "/")
	ctx := context.Background()

	ref := fmt.Sprintf("%s/%s", repositoryUrl, chartName)

	// Crée un nouveau référentiel distant
	repo, err := remote.NewRepository(ref)
	if err != nil {
		logrus.Errorf("Erreur lors de la création du référentiel distant : %s\n%v", ref, err)
		return nil, err
	}

	var tags []string
	// Récupère la liste des tags
	err = repo.Tags(ctx, "", func(tagsFound []string) error {
		// tags = filterVersions(tagsFound)
		tags = tagsFound
		return nil
	})

	if err != nil {
		logrus.Tracef("Erreur lors de la récupération des tags : %v", err)
		return nil, err
	}
	return tags, nil
}
