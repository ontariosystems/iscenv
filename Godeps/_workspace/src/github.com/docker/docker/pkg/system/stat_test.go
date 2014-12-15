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
	"syscall"
	"testing"
)

func TestFromStatT(t *testing.T) {
	file, _, _, dir := prepareFiles(t)
	defer os.RemoveAll(dir)

	stat := &syscall.Stat_t{}
	err := syscall.Lstat(file, stat)

	s, err := fromStatT(stat)
	if err != nil {
		t.Fatal(err)
	}

	if stat.Mode != s.Mode() {
		t.Fatal("got invalid mode")
	}
	if stat.Uid != s.Uid() {
		t.Fatal("got invalid uid")
	}
	if stat.Gid != s.Gid() {
		t.Fatal("got invalid gid")
	}
	if stat.Rdev != s.Rdev() {
		t.Fatal("got invalid rdev")
	}
	if stat.Mtim != s.Mtim() {
		t.Fatal("got invalid mtim")
	}
}
