package versions

import (
	"github.com/sgreben/versions/pkg/simplegit"
	"io"
)

type ThingPartsSourceGit struct {
	Repository simplegit.Repository
}

func (t *ThingPartsSourceGit) Fetch(fileName string) (io.Reader, error) {
	raw, err := t.Repository.Raw()
	if err != nil {
		return nil, err
	}
	wt, err := raw.Worktree()
	if err != nil {
		return nil, err
	}
	f, err := wt.Filesystem.Open(fileName)
	if err != nil {
		return nil, err
	}
	return f, nil
}
