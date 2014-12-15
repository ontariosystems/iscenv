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
	"os"
	"testing"
)

func TestLstat(t *testing.T) {
	file, invalid, _, dir := prepareFiles(t)
	defer os.RemoveAll(dir)

	statFile, err := Lstat(file)
	if err != nil {
		t.Fatal(err)
	}
	if statFile == nil {
		t.Fatal("returned empty stat for existing file")
	}

	statInvalid, err := Lstat(invalid)
	if err == nil {
		t.Fatal("did not return error for non-existing file")
	}
	if statInvalid != nil {
		t.Fatal("returned non-nil stat for non-existing file")
	}
}
