package versions

import (
	"github.com/sgreben/versions/pkg/semver"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

func FetchFromGitTags(url string) (VersionsWithSources, error) {

	s := memory.NewStorage()
	r, err := git.Clone(s, nil, &git.CloneOptions{
		URL:        url,
		NoCheckout: true,
		Tags:       git.AllTags,
	})
	if err != nil {
		return nil, err
	}
	iter, err := r.Tags()
	if err != nil {
		return nil, err
	}

	out := VersionsWithSources{}

	iter.ForEach(func(tag *plumbing.Reference) (errOut error) {
		version, err := semver.Parse(tag.Name().Short())
		if err != nil {
			return
		}
		source := VersionSourceGit{
			RepositoryURL: url,
			Reference:     tag.Name().String(),
		}
		out = append(out, VersionWithSource{
			Version: version,
			Source: VersionSource{
				Git: &source,
			},
		})
		return
	})
	return out, nil
}
