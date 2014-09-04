/*
Copyright 2014 Ontario Systems

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

package version

import (
	"sort"
)

// Sorts a string slice of version number strings using version.CompareSimple()
//
// Example:
//     version.Sort([]string{"1.10-dev", "1.0rc1", "1.0", "1.0-dev"})
//     Returns []string{"1.0-dev", "1.0rc1", "1.0", "1.10-dev"}
//
func Sort(versionStrings []string) {
	versions := versionSlice(versionStrings)
	sort.Sort(versions)
}

type versionSlice []string

func (s versionSlice) Len() int {
	return len(s)
}

func (s versionSlice) Less(i, j int) bool {
	cmp := CompareSimple(s[i], s[j])
	if cmp == 0 {
		return s[i] < s[j]
	}
	return cmp < 0
}

func (s versionSlice) Swap(i, j int) {
	tmp := s[j]
	s[j] = s[i]
	s[i] = tmp
}
