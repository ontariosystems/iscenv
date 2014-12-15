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

package system

import (
	"strings"
	"testing"

	"github.com/docker/docker/pkg/units"
)

func TestMemInfo(t *testing.T) {
	const input = `
	MemTotal:      1 kB
	MemFree:       2 kB
	SwapTotal:     3 kB
	SwapFree:      4 kB
	Malformed1:
	Malformed2:    1
	Malformed3:    2 MB
	Malformed4:    X kB
	`
	meminfo, err := parseMemInfo(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if meminfo.MemTotal != 1*units.KiB {
		t.Fatalf("Unexpected MemTotal: %d", meminfo.MemTotal)
	}
	if meminfo.MemFree != 2*units.KiB {
		t.Fatalf("Unexpected MemFree: %d", meminfo.MemFree)
	}
	if meminfo.SwapTotal != 3*units.KiB {
		t.Fatalf("Unexpected SwapTotal: %d", meminfo.SwapTotal)
	}
	if meminfo.SwapFree != 4*units.KiB {
		t.Fatalf("Unexpected SwapFree: %d", meminfo.SwapFree)
	}
}
