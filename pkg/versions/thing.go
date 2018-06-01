package versions

import (
	"errors"
	"io"

	"github.com/sgreben/versions/pkg/semver"
)

// A ThingName is a *unique* identifier for a thing
type ThingName string

// A thing has a name, a set of versions, and may want other things
type Thing struct {
	Name     ThingName
	Versions ThingVersions
	Wants    ThingWants
	Parts    ThingParts
}

// ThingParts is the stuff that makes up the thing
type ThingParts struct {
	Cached *func(string) (io.Reader, error)
	ThingPartsSource
}

func (t ThingParts) FetchOrCached(part string) (io.Reader, error) {
	if t.Cached != nil {
		return (*t.Cached)(part)
	}
	return t.ThingPartsSource.Fetch(part)
}

// ThingPartsSource specifies how to obtain a thing's parts
type ThingPartsSource struct {
	FromGit *ThingPartsSourceGit
}

func (t ThingPartsSource) Fetch(part string) (io.Reader, error) {
	switch {
	case t.FromGit != nil:
		return t.FromGit.Fetch(part)
	default:
		return nil, errors.New("no ThingPartsSource defined")
	}
}

type ThingWants struct {
	Cached *[]WantedThing
	ThingWantsSource
}

func (t *ThingWants) CachedOrFetch(thing Thing) ([]WantedThing, error) {
	if t.Cached != nil {
		return *t.Cached, nil
	}
	return t.ThingWantsSource.Fetch(thing)
}

// ThingWantsSource specifies how to determine which other things a thing wants
type ThingWantsSource struct {
	Part           string
	FromSerialized *ThingWantsSourceSerialized
}

func (t *ThingWantsSource) Fetch(thing Thing) ([]WantedThing, error) {
	part, err := thing.Parts.FetchOrCached(t.Part)
	if err != nil {
		return nil, err
	}
	switch {
	case t.FromSerialized != nil:
		return t.FromSerialized.Fetch(part)
	default:
		return nil, errors.New("no ThingWantsSource defined")
	}
}

type WantedThing struct {
	Thing
	WantedVersion *semver.Constraints
}

type ThingVersions struct {
	Cached *VersionsWithSources
	ThingVersionsSource
}

func (t *ThingVersions) CachedOrFetch() (VersionsWithSources, error) {
	if t.Cached != nil {
		return *t.Cached, nil
	}
	svs, err := t.ThingVersionsSource.Fetch()
	if err != nil {
		return nil, err
	}
	t.Cached = &svs
	return svs, nil
}

type ThingVersionsSource struct {
	FromGit    *ThingVersionsSourceGit
	FromDocker *ThingVersionsSourceDocker
}

func (t ThingVersionsSource) Fetch() (VersionsWithSources, error) {
	switch {
	case t.FromGit != nil:
		return t.FromGit.Fetch()
	case t.FromDocker != nil:
		return t.FromDocker.Fetch()
	default:
		return nil, errors.New("no ThingVersionsSource defined")
	}
}
