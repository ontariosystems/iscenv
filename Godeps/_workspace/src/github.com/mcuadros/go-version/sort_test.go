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
	"reflect"
	"testing"
)

func TestSort(t *testing.T) {
	testcases := []struct {
		input  []string
		output []string
	}{
		{
			input: []string{
				"Package-0.4.tar.gz",
				"Package-0.1.tar.gz",
				"Package-0.10.1.tar.gz",
				"Package-0.10.tar.gz",
				"Package-0.2.tar.gz",
				"Package-0.3.1.tar.gz",
				"Package-0.3.2.tar.gz",
				"Package-0.3.tar.gz",
			},
			output: []string{
				"Package-0.1.tar.gz",
				"Package-0.2.tar.gz",
				"Package-0.3.tar.gz",
				"Package-0.3.1.tar.gz",
				"Package-0.3.2.tar.gz",
				"Package-0.4.tar.gz",
				"Package-0.10.tar.gz",
				"Package-0.10.1.tar.gz",
			},
		},
		{
			input: []string{
				"1.0-dev",
				"1.0a1",
				"1.0b1",
				"1.0RC1",
				"1.0rc1",
				"1.0",
				"1.0pl1",
				"1.1-dev",
				"1.2",
				"1.10",
			},
			output: []string{
				"1.0-dev",
				"1.0a1",
				"1.0b1",
				"1.0RC1",
				"1.0rc1",
				"1.0",
				"1.0pl1",
				"1.1-dev",
				"1.2",
				"1.10",
			},
		},
	}

	for _, testcase := range testcases {
		Sort(testcase.input)
		if !reflect.DeepEqual(testcase.input, testcase.output) {
			t.Errorf("Expected output %+v did not match actual %+v", testcase.output, testcase.input)
		}
	}
}
