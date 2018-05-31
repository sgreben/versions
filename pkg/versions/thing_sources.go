package versions

import (
	"io"
	"log"

	"gopkg.in/src-d/go-billy.v4"

	"gopkg.in/src-d/go-billy.v4/memfs"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

type ThingSourceGit struct {
	RepositoryURL string
	Cached        *billy.Filesystem
}

func (t *ThingSourceGit) Fetch(fileName string) (io.Reader, error) {
	if t.Cached == nil {
		s := memory.NewStorage()
		fs := memfs.New()
		r, err := git.Clone(s, fs, &git.CloneOptions{
			URL:   t.RepositoryURL,
			Tags:  git.AllTags,
			Depth: 1,
		})
		if err != nil {
			return nil, err
		}
		w, err := r.Worktree()
		if err != nil {
			return nil, err
		}
		t.Cached = &w.Filesystem
	}
	fs := *t.Cached
	log.Println(t.RepositoryURL, "file:", fileName)
	f, err := fs.Open(fileName)
	if err != nil {
		return nil, err
	}
	return f, nil
}
