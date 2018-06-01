package versions

import "fmt"

// http://archive.is/aDhCk#selection-135.1-161.208

type MVSBuildList struct {
	Versions map[ThingName]VersionWithSource
	Things   []Thing
}

func (mvs MVSBuildList) Construct() (errs []error) {
	if mvs.Versions == nil {
		mvs.Versions = map[ThingName]VersionWithSource{}
	}
	oldThings := make([]Thing, len(mvs.Things))
	copy(oldThings, mvs.Things)
	mvs.Things = mvs.Things[:0]
	for _, oldThing := range oldThings {
		mvs.Things = append(mvs.Things, oldThing)
		wants, err := oldThing.Wants.CachedOrFetch(oldThing)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		for _, want := range wants {
			mvs.Things = append(mvs.Things, want.Thing)
			versions, err := want.Thing.Versions.CachedOrFetch()
			if err != nil {
				errs = append(errs, err)
				continue
			}
			latestMatchingVersion := versions.LatestMatching(want.WantedVersion)
			if latestMatchingVersion == nil {
				errs = append(errs, fmt.Errorf("no matching version for %s (%v)", want.Thing.Name, want.WantedVersion))
				continue
			}
			if currentVersion, ok := mvs.Versions[want.Name]; ok {
				if currentVersion.Version.LessThan(latestMatchingVersion.Version) {
					mvs.Versions[want.Name] = *latestMatchingVersion
				}
			}
		}
	}
	return
}
