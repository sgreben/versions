package versions

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/BurntSushi/toml"

	"github.com/sgreben/versions/pkg/semver"

	"github.com/sgreben/versions/pkg/jsonpath"
)

type ThingWantsSourceSerialized struct {
	Deserializer
	ExtractPath *jsonpath.JSONPath
}

type Deserializer struct {
	FromJSON *JSONDeserializer
	FromTOML *TOMLDeserializer
}

func (d Deserializer) Deserialize(r io.Reader) (interface{}, error) {
	switch {
	case d.FromJSON != nil:
		return d.FromJSON.Deserialize(r)
	case d.FromTOML != nil:
		return d.FromTOML.Deserialize(r)
	default:
		return nil, errors.New("no Deserializer defined")
	}
}

type JSONDeserializer struct{}

func (d JSONDeserializer) Deserialize(r io.Reader) (data interface{}, err error) {
	err = json.NewDecoder(r).Decode(&data)
	return
}

type TOMLDeserializer struct{}

func (d TOMLDeserializer) Deserialize(r io.Reader) (data interface{}, err error) {
	_, err = toml.DecodeReader(r, &data)
	return
}

func (t *ThingWantsSourceSerialized) Fetch(r io.Reader) ([]WantedThing, error) {
	var out []WantedThing
	data, err := t.Deserialize(r)
	if err != nil {
		return nil, err
	}
	results, err := t.ExtractPath.AllowMissingKeys(true).FindResults(data)
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
