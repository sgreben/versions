package versions

import (
	"sort"

	"github.com/sgreben/versions/pkg/semver"
)

type VersionWithSource struct {
	Version *semver.Version
	Source  VersionSource
}

type VersionSource struct {
	Git    *VersionSourceGit    `json:",omitempty"`
	Docker *VersionSourceDocker `json:",omitempty"`
}

type VersionsWithSources []VersionWithSource

func (versions VersionsWithSources) LatestMatching(constraints *semver.Constraints) *VersionWithSource {
	sort.Sort(versions)
	for i := 0; i < len(versions); i++ {
		candidate := versions[len(versions)-1-i]
		if constraints.Check(candidate.Version) {
			return &candidate
		}
	}
	return nil
}

func (c VersionsWithSources) Len() int {
	return len(c)
}

func (c VersionsWithSources) Less(i, j int) bool {
	return c[i].Version.LessThan(c[j].Version)
}

func (c VersionsWithSources) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
