package main

import (
	"log"
	"os"
	"sort"

	"github.com/sgreben/versions/pkg/versions"
)

func fetchFromGitCmd(url string, limit int) {
	svs, err := versions.FetchFromGitTags(url)
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

func fetchFromDockerCmd(repository string, limit int) {
	svs, err := versions.FetchFromDocker(repository)
	if err != nil {
		log.Println(err)
	}
	sort.Sort(svs)
	if limit > 0 && len(svs) > limit {
		svs = svs[len(svs)-limit:]
	}
	err = jsonEncode(svs, os.Stdout)
	if err != nil {
		log.Println(err)
	}
}
