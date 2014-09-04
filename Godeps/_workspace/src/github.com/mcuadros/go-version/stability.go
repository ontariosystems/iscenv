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
	"regexp"
	"strings"
)

const (
	Development = iota
	Alpha
	Beta
	RC
	Stable
)

func expandStability(stability string) string {
	stability = strings.ToLower(stability)

	switch stability {
	case "a":
		return "alpha"
	case "b":
		return "beta"
	case "p":
		return "patch"
	case "pl":
		return "patch"
	case "rc":
		return "RC"
	}

	return stability
}

func parseStability(version string) string {
	version = regexp.MustCompile(`(?i)#.+$`).ReplaceAllString(version, " ")
	version = strings.ToLower(version)

	if strings.HasPrefix(version, "dev-") || strings.HasSuffix(version, "-dev") {
		return "dev"
	}

	result := RegFind(`(?i)^v?(\d{1,3})(\.\d+)?(\.\d+)?(\.\d+)?`+modifierRegex+`$`, version)
	if result != nil {
		if len(result) > 3 {
			return "dev"
		}
	}

	if result[1] != "" {
		if "beta" == result[1] || "b" == result[1] {
			return "beta"
		}
		if "alpha" == result[1] || "a" == result[1] {
			return "alpha"
		}
		if "rc" == result[1] {
			return "RC"
		}
	}

	return "stable"
}

func GetStability(version string) int {
	result := RegFind(`(?i)(stable|RC|beta|alpha|dev)`, Normalize(version))
	if len(result) == 0 {
		return Stable
	}

	switch result[1] {
	case "dev":
		return Development
	case "alpha":
		return Alpha
	case "beta":
		return Beta
	case "RC":
		return RC
	}

	return Stable
}
