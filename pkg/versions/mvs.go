package versions

import (
	"fmt"

	"github.com/sgreben/versions/pkg/semver"
)

// ConstraintsForName maps a "package" name to its version constraint
type ConstraintsForName map[string]*semver.Constraints

// ConstraintsForVersion maps a package version to its set of dependency constraints
type ConstraintsForVersion map[string]ConstraintsForName

// ConstraintGraph maps a package name to its set of versions
type ConstraintGraph map[string]ConstraintsForVersion

// SelectMVS uses the MVS (https://research.swtch.com/vgo-mvs) algorithm to solve a version constraint graph
func (d ConstraintGraph) SelectMVS(forName string, forVersion *semver.Version) (map[string]*semver.Version, error) {
	work := []struct {
		Name    string
		Version *semver.Version
	}{{
		Name:    forName,
		Version: forVersion,
	}}
	out := map[string]*semver.Version{}
	for len(work) > 0 {
		item := work[0]
		forThing, ok := d[item.Name][item.Version.String()]
		if !ok {
			forThing, ok = d[item.Name][item.Version.Original]
		}
		work = work[1:]

		if !ok {
			return nil, fmt.Errorf("no version of %q matching constraint %v", item.Name, item.Version)
		}
		for name, wanted := range forThing {
			var available semver.Collection
			for k := range d[name] {
				v, _ := semver.Parse(string(k))
				available = append(available, v)
			}
			matching := wanted.OldestMatching(available)
			if matching == nil {
				if item.Name == forName {
					return nil, fmt.Errorf("no version of %q matching constraint %q of %q (available: %v)", name, wanted, item.Name, available)
				}
				return nil, fmt.Errorf("no version of %q matching constraint %q of %q:%q", name, wanted, item.Name, item.Version.String())
			}
			if current, ok := out[name]; ok {
				if current.LessThan(matching) {
					out[name] = matching
					itemCopy := item
					itemCopy.Name = name
					itemCopy.Version = matching
					work = append(work, itemCopy)
				}
				continue
			}
			out[name] = matching
			itemCopy := item
			itemCopy.Name = name
			itemCopy.Version = matching
			work = append(work, itemCopy)
		}
	}
	return out, nil
}
