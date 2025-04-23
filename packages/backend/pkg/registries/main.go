package registries

import (
	"argocd-watcher/pkg/cachex"
	"context"

	"github.com/blang/semver/v4"
)

func GetGreaterTags(registry, chartName string, minVersion semver.Version) ([]string, error) {
	tags := []string{}

	chart, err := Search(registry, chartName)

	if err != nil {
		return nil, err
	}
	for _, version := range chart.Tags {
		candidateVersion, _ := semver.Parse(version)
		if minVersion.LT(candidateVersion) {
			tags = append(tags, version)
		}
	}
	return tags, nil
}

var ChartCache *cachex.CacheService[Chart]

func Init() {
	ctx := context.Background()
	ChartCache = cachex.NewCacheService[Chart](ctx, "cache:invalidation")
}
