package main

import (
	"github.com/sgreben/versions/pkg/simpledocker"
	"github.com/sgreben/versions/pkg/simplegit"
	"github.com/sgreben/versions/pkg/versions"
	git "gopkg.in/src-d/go-git.v4"
	"log"
	"os"
	"sort"
)

func fetchFromGit(url string, limit int) (versions.VersionsWithSources, error) {
	vs := versions.ThingVersionsSourceGit{
		Repository: simplegit.Repository{
			URL: url,
			CloneOptions: &git.CloneOptions{
				NoCheckout: true,
				Depth:      0,
				Tags:       git.AllTags,
			},
		},
	}
	svs, err := vs.Fetch()
	if err != nil {
		return nil, err
	}
	sort.Sort(svs)
	if limit > 0 && len(svs) > limit {
		svs = svs[len(svs)-limit:]
	}
	return svs, nil
}

func fetchFromGitCmd(url string, limit int) {
	svs, err := fetchFromGit(url, limit)
	if err != nil {
		log.Println(err)
		exit.NonzeroBecause = append(exit.NonzeroBecause, err.Error())
	}
	sort.Sort(svs)
	if limit > 0 && len(svs) > limit {
		svs = svs[len(svs)-limit:]
	}
	err = jsonEncode(svs, os.Stdout)
	if err != nil {
		log.Println(err)
		exit.NonzeroBecause = append(exit.NonzeroBecause, err.Error())
	}
}

func fetchFromDocker(repository string, limit int) (versions.VersionsWithSources, error) {
	vs := versions.ThingVersionsSourceDocker{
		Repository: &simpledocker.Repository{
			URL: repository,
		},
	}
	svs, err := vs.Fetch()
	if err != nil {
		return nil, err
	}
	sort.Sort(svs)
	if limit > 0 && len(svs) > limit {
		svs = svs[len(svs)-limit:]
	}
	return svs, err
}

func fetchFromDockerCmd(repository string, limit int) {
	svs, err := fetchFromDocker(repository, limit)
	if err != nil {
		log.Println(err)
	}
	err = jsonEncode(svs, os.Stdout)
	if err != nil {
		log.Println(err)
	}
}
