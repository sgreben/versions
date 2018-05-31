package versions

type VersionSourceGit struct {
	RepositoryURL string
	Reference     string
}

type VersionSourceDocker struct {
	Image string
	Tag   string
}
