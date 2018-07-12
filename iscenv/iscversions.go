/*
Copyright 2016 Ontario Systems

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package iscenv

import (
	"sort"

	"github.com/mcuadros/go-version"
)

// ISCVersions is a slice of ISCVersion
type ISCVersions []*ISCVersion

// Len returns the number of versions
func (evs ISCVersions) Len() int { return len(evs) }

// Swap interchanges two versions in the list
func (evs ISCVersions) Swap(i, j int) { evs[i], evs[j] = evs[j], evs[i] }

// Less returns whether version in one index is less than the version in the other index
func (evs ISCVersions) Less(i, j int) bool {
	cmp := version.CompareSimple(version.Normalize(evs[i].Version), version.Normalize(evs[j].Version))
	if cmp == 0 {
		return evs[i].Version < evs[j].Version
	}

	return cmp < 0
}

// AddIfMissing adds a version to the list if it isn't already included
func (evs *ISCVersions) AddIfMissing(ev *ISCVersion) bool {
	if !evs.Exists(ev.Version) {
		*evs = append(*evs, ev)
		return true
	}

	return false
}

// Latest finds and returns the last version in the list
func (evs ISCVersions) Latest() *ISCVersion {
	return evs[len(evs)-1]
}

// Exists returns whether the provided version exists in the list
func (evs ISCVersions) Exists(versionString string) bool {
	return evs.Find(versionString) != nil
}

// Find finds and returns the ISCVersion for the provided version string
func (evs ISCVersions) Find(versionString string) *ISCVersion {
	for _, version := range evs {
		if version.Version == versionString {
			return version
		}
	}

	return nil
}

// Sort sorts the versions in the list
func (evs *ISCVersions) Sort() {
	sort.Sort(*evs)
}
