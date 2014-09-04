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

import "fmt"

type ChangeType int

const (
	ChangeModify ChangeType = iota
	ChangeAdd
	ChangeDelete
)

// Change represents a change in a container.
//
// See http://goo.gl/DpGyzK for more details.
type Change struct {
	Path string
	Kind ChangeType
}

func (change *Change) String() string {
	var kind string
	switch change.Kind {
	case ChangeModify:
		kind = "C"
	case ChangeAdd:
		kind = "A"
	case ChangeDelete:
		kind = "D"
	}
	return fmt.Sprintf("%s %s", kind, change.Path)
}
