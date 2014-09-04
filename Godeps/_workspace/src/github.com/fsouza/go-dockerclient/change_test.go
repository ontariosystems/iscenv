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

// Copyright 2014 go-dockerclient authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package docker

import (
	"testing"
)

func TestChangeString(t *testing.T) {
	var tests = []struct {
		change   Change
		expected string
	}{
		{Change{"/etc/passwd", ChangeModify}, "C /etc/passwd"},
		{Change{"/etc/passwd", ChangeAdd}, "A /etc/passwd"},
		{Change{"/etc/passwd", ChangeDelete}, "D /etc/passwd"},
		{Change{"/etc/passwd", 33}, " /etc/passwd"},
	}
	for _, tt := range tests {
		if got := tt.change.String(); got != tt.expected {
			t.Errorf("Change.String(): want %q. Got %q.", tt.expected, got)
		}
	}
}
