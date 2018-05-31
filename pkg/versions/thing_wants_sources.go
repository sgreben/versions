package versions

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/BurntSushi/toml"

	"github.com/sgreben/versions/pkg/semver"

	"github.com/sgreben/versions/pkg/jsonpath"
)

type ThingWantsSourceTOML struct {
	JSONPath *jsonpath.JSONPath
}

func (t *ThingWantsSourceTOML) Fetch(r io.Reader) ([]WantedThing, error) {
	var out []WantedThing
	var data interface{}
	_, err := toml.DecodeReader(r, &data)
	if err != nil {
		return nil, err
	}
	results, err := t.JSONPath.AllowMissingKeys(true).FindResults(data)
	if err != nil {
		return nil, err
	}
	for _, result := range results {
		if len(result) < 1 {
			return nil, fmt.Errorf("invalid requirement: %v", result)
		}
		wantedVersionString := "*"
		if len(result) >= 2 {
			wantedVersionValue := result[1]
			if !wantedVersionValue.IsNil() {
				wantedVersionString = fmt.Sprint(wantedVersionValue.Interface())
			}
		}
		name := ThingName(fmt.Sprint(result[0].Interface()))
		wantedVersion, err := semver.ParseConstraint(wantedVersionString)
		if err != nil {
			return nil, err
		}
		out = append(out, WantedThing{
			Thing:         Thing{Name: name},
			WantedVersion: wantedVersion,
		})
	}
	return out, nil
}

type ThingWantsSourceJSON struct {
	JSONPath *jsonpath.JSONPath
}

func (t *ThingWantsSourceJSON) Fetch(r io.Reader) ([]WantedThing, error) {
	var out []WantedThing
	dec := json.NewDecoder(r)
	var data interface{}
	err := dec.Decode(&data)
	if err != nil {
		return nil, err
	}
	t.JSONPath.AllowMissingKeys(true)
	results, err := t.JSONPath.AllowMissingKeys(true).FindResults(data)
	if err != nil {
		return nil, err
	}
	for _, result := range results {
		if len(result) < 2 {
			return nil, fmt.Errorf("invalid requirement: %v", result)
		}
		name := ThingName(result[0].String())
		wantedVersionString := result[1].String()
		wantedVersion, err := semver.ParseConstraint(wantedVersionString)
		if err != nil {
			return nil, err
		}
		out = append(out, WantedThing{
			Thing:         Thing{Name: name},
			WantedVersion: wantedVersion,
		})
	}
	return out, nil
}
