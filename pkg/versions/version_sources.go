package versions

import (
	"github.com/sgreben/versions/pkg/simplegit"
)

type VersionSourceGit struct {
	Repository simplegit.Repository
	Reference  string
}

type VersionSourceDocker struct {
	Image string
	Tag   string
}
