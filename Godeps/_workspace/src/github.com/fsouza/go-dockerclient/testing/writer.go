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

package testing

import (
	"encoding/binary"
	"errors"
	"io"
)

type stdType [8]byte

var (
	stdin  = stdType{0: 0}
	stdout = stdType{0: 1}
	stderr = stdType{0: 2}
)

type stdWriter struct {
	io.Writer
	prefix  stdType
	sizeBuf []byte
}

func (w *stdWriter) Write(buf []byte) (n int, err error) {
	if w == nil || w.Writer == nil {
		return 0, errors.New("Writer not instanciated")
	}
	binary.BigEndian.PutUint32(w.prefix[4:], uint32(len(buf)))
	buf = append(w.prefix[:], buf...)

	n, err = w.Writer.Write(buf)
	return n - 8, err
}

func newStdWriter(w io.Writer, t stdType) *stdWriter {
	if len(t) != 8 {
		return nil
	}
	return &stdWriter{Writer: w, prefix: t, sizeBuf: make([]byte, 4)}
}
