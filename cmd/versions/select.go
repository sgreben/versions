package main

import (
	"fmt"
	"os"

	"github.com/sgreben/versions/pkg/semver"
	"github.com/sgreben/versions/pkg/versions"
)

func selectSingleCmd(constraint string, versions []string) {
	c, err := semver.ParseConstraint(constraint)
	if err != nil {
		exit.NonzeroBecause = append(exit.NonzeroBecause, fmt.Sprintf("cannot parse constraint %q: %v", constraint, err))
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

func selectAllCmd(constraint string, versions []string) {
	c, err := semver.ParseConstraint(constraint)
	if err != nil {
		exit.NonzeroBecause = append(exit.NonzeroBecause, fmt.Sprintf("cannot parse constraint %q: %v", constraint, err))
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
	solution := c.AllMatching(svs)
	jsonEncode(solution, os.Stdout)
}

func selectMvsCmd(targetName string, g versions.ConstraintGraph) {
	for forVersionString := range g[targetName] {
		forVersion, err := semver.Parse(forVersionString)
		if err != nil {
			forVersion = &semver.Version{}
		}
		mvsOutput, err := g.SelectMVS(targetName, forVersion)
		if err != nil {
			exit.NonzeroBecause = append(exit.NonzeroBecause, fmt.Sprintf("mvs failed: %v", err))
			return
		}
		jsonEncode(mvsOutput, os.Stdout)
		return
	}
}
