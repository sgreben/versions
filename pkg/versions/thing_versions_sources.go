package versions

type ThingVersionsSourceGit struct {
	RepositoryURL string
}

func (t *ThingVersionsSourceGit) Fetch() (VersionsWithSources, error) {
	return FetchFromGitTags(t.RepositoryURL)
}

type ThingVersionsSourceDocker struct {
	RepositoryURL string
}

func (t *ThingVersionsSourceDocker) Fetch() (VersionsWithSources, error) {
	return FetchFromDocker(t.RepositoryURL)
}
