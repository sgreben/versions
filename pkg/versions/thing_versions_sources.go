package versions

import (
	"github.com/sgreben/versions/pkg/semver"
	"github.com/sgreben/versions/pkg/simpledocker"
	"github.com/sgreben/versions/pkg/simplegit"
)

type ThingVersionsSourceGit struct {
	Repository simplegit.Repository
}

func (t ThingVersionsSourceGit) Fetch() (VersionsWithSources, error) {
	tags, err := t.Repository.Tags()
	if err != nil {
		return nil, err
	}
	out := make(VersionsWithSources, 0, len(tags))
	for _, tag := range tags {
		version, err := semver.Parse(tag.Name)
		if err != nil {
			continue
		}
		source := VersionSourceGit{
			Repository: t.Repository,
			Reference:  tag.Reference,
		}
		out = append(out, VersionWithSource{
			Version: version,
			Source: VersionSource{
				Git: &source,
			},
		})
	}
	return out, nil
}

type ThingVersionsSourceDocker struct {
	Repository *simpledocker.Repository
}

func (t ThingVersionsSourceDocker) Fetch() (VersionsWithSources, error) {
	tags, err := t.Repository.Tags()
	if err != nil {
		return nil, err
	}
	out := make(VersionsWithSources, 0, len(tags))
	for _, tag := range tags {
		version, err := semver.Parse(tag.Name)
		if err != nil {
			continue
		}
		source := VersionSourceDocker{
			Tag:   tag.Name,
			Image: tag.Image,
		}
		out = append(out, VersionWithSource{
			Version: version,
			Source: VersionSource{
				Docker: &source,
			},
		})
	}
	return out, nil
}
