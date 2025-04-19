package registries

import "github.com/blang/semver/v4"

func GetGreaterTags(index *IndexFile, registry, chart string, minVersion semver.Version) []string {
	tags := []string{}

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

func GetTags(registry, chart string) []string {
	tags := []string{}
	index, _ := Search(registry)
	if entry, ok := index.Entries[chart]; ok {
		for _, version := range entry {
			tags = append(tags, version.Version)
		}
	}
	return tags
}
