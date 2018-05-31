package versions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/sgreben/versions/pkg/semver"
)

func FetchFromDocker(repository string) (VersionsWithSources, error) {
	i := strings.IndexByte(repository, '/')
	n := strings.Count(repository, "/")
	switch {
	case n == 0:
		user := "library"
		image := repository
		return FetchFromDockerhubTags(user, image)
	case n == 1:
		user := repository[:i]
		image := repository[i+1:]
		return FetchFromDockerhubTags(user, image)
	case n >= 2:
		registry, image := repository[:i], repository[i+1:]
		return FetchFromDockerTags(registry, image)
	default:
		return nil, fmt.Errorf("cannot determine registry: %s", repository)
	}
}

func FetchFromDockerhubTags(user, image string) (VersionsWithSources, error) {
	url := fmt.Sprintf("https://registry.hub.docker.com/v2/repositories/%s/%s/tags/", user, image)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	out := VersionsWithSources{}
	type response struct {
		Count   int     `json:"count"`
		Next    *string `json:"next"`
		Results []struct {
			Name string `json:"name"`
		} `json:"results"`
	}
	dec := json.NewDecoder(resp.Body)
	var apiResp response
	err = dec.Decode(&apiResp)
	if err != nil {
		return nil, err
	}
	for {
		for _, result := range apiResp.Results {
			tag := result.Name
			version, err := semver.Parse(tag)
			if err != nil {
				continue
			}
			source := VersionSourceDocker{
				Image: fmt.Sprintf("%s/%s:%s", user, image, tag),
				Tag:   tag,
			}
			out = append(out, VersionWithSource{
				Version: version,
				Source: VersionSource{
					Docker: &source,
				},
			})
		}
		if apiResp.Next == nil {
			break
		}
		url = *apiResp.Next
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		dec = json.NewDecoder(resp.Body)
		apiResp = response{}
		err = dec.Decode(&apiResp)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func FetchFromDockerTags(registry, image string) (VersionsWithSources, error) {
	url := fmt.Sprintf("https://%s/v2/%s/tags/list", registry, image)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	out := VersionsWithSources{}
	type response struct {
		Tags []string `json:"tags"`
	}
	dec := json.NewDecoder(resp.Body)
	var apiResp response
	err = dec.Decode(&apiResp)
	if err != nil {
		return nil, err
	}
	for _, tag := range apiResp.Tags {
		version, err := semver.Parse(tag)
		if err != nil {
			continue
		}
		source := VersionSourceDocker{
			Image: fmt.Sprintf("%s/%s:%s", registry, image, tag),
			Tag:   tag,
		}
		out = append(out, VersionWithSource{
			Version: version,
			Source: VersionSource{
				Docker: &source,
			},
		})
	}
	return out, nil
}
