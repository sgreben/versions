package main

import (
	"fmt"
	"os"

	"github.com/sgreben/versions/pkg/semver"
)

func selectSingleCmd(constraint string, versions []string) {
	c, err := semver.ParseConstraint(constraint)
	if err != nil {
		exit.NonzeroBecause = append(exit.NonzeroBecause, "no matching version")
		return
	}
	svs := make(semver.Collection, 0, len(versions))
	for _, v := range versions {
		sv, err := semver.Parse(v)
		if err != nil {
			exit.NonzeroBecause = append(exit.NonzeroBecause, fmt.Sprintf(`"%s": %v`, v, err))
			continue
		}
		svs = append(svs, sv)
	}
	solution := c.LatestMatching(svs)
	if solution == nil {
		exit.NonzeroBecause = append(exit.NonzeroBecause, "no matching version")
		return
	}
	jsonEncode(solution.String(), os.Stdout)
}

func selectMvsCmd(_ string) {
	exit.NonzeroBecause = append(exit.NonzeroBecause, "not implemented")
}
