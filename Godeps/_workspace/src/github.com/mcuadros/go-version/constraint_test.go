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
	"testing"
)

func TestGetOperator(t *testing.T) {
	constraint := NewConstrain("=", "1.0.0")
	out := "="

	if x := constraint.GetOperator(); x != "=" {
		t.Errorf("FAIL: GetOperator() = {%s}: want {%s}", x, out)
	}
}

func TestGetVersion(t *testing.T) {
	constraint := NewConstrain("=", "1.0.0")
	out := "1.0.0"

	if x := constraint.GetVersion(); x != "1.0.0" {
		t.Errorf("FAIL: GetVersion() = {%s}: want {%s}", x, out)
	}
}

func TestString(t *testing.T) {
	constraint := NewConstrain("=", "1.0.0")
	out := "= 1.0.0"

	if x := constraint.String(); x != out {
		t.Errorf("FAIL: String() = {%s}: want {%s}", x, out)
	}
}

func TestMatchSuccess(t *testing.T) {
	constraint := NewConstrain("=", "1.0.0")
	out := true

	if x := constraint.Match("1.0"); x != out {
		t.Errorf("FAIL: Match() = {%s}: want {%s}", x, out)
	}
}

func TestMatchFail(t *testing.T) {
	constraint := NewConstrain("=", "1.0.0")
	out := false

	if x := constraint.Match("2.0"); x != out {
		t.Errorf("FAIL: Match() = {%s}: want {%s}", x, out)
	}
}
