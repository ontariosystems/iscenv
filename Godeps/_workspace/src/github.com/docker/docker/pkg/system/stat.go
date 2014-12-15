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
	"syscall"
)

type Stat struct {
	mode uint32
	uid  uint32
	gid  uint32
	rdev uint64
	size int64
	mtim syscall.Timespec
}

func (s Stat) Mode() uint32 {
	return s.mode
}

func (s Stat) Uid() uint32 {
	return s.uid
}

func (s Stat) Gid() uint32 {
	return s.gid
}

func (s Stat) Rdev() uint64 {
	return s.rdev
}

func (s Stat) Size() int64 {
	return s.size
}

func (s Stat) Mtim() syscall.Timespec {
	return s.mtim
}

func (s Stat) GetLastModification() syscall.Timespec {
	return s.Mtim()
}
