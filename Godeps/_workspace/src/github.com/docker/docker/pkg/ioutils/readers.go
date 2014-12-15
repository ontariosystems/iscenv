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

package ioutils

import (
	"bytes"
	"io"
	"sync"
)

type readCloserWrapper struct {
	io.Reader
	closer func() error
}

func (r *readCloserWrapper) Close() error {
	return r.closer()
}

func NewReadCloserWrapper(r io.Reader, closer func() error) io.ReadCloser {
	return &readCloserWrapper{
		Reader: r,
		closer: closer,
	}
}

type readerErrWrapper struct {
	reader io.Reader
	closer func()
}

func (r *readerErrWrapper) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	if err != nil {
		r.closer()
	}
	return n, err
}

func NewReaderErrWrapper(r io.Reader, closer func()) io.Reader {
	return &readerErrWrapper{
		reader: r,
		closer: closer,
	}
}

type bufReader struct {
	sync.Mutex
	buf      *bytes.Buffer
	reader   io.Reader
	err      error
	wait     sync.Cond
	drainBuf []byte
}

func NewBufReader(r io.Reader) *bufReader {
	reader := &bufReader{
		buf:      &bytes.Buffer{},
		drainBuf: make([]byte, 1024),
		reader:   r,
	}
	reader.wait.L = &reader.Mutex
	go reader.drain()
	return reader
}

func NewBufReaderWithDrainbufAndBuffer(r io.Reader, drainBuffer []byte, buffer *bytes.Buffer) *bufReader {
	reader := &bufReader{
		buf:      buffer,
		drainBuf: drainBuffer,
		reader:   r,
	}
	reader.wait.L = &reader.Mutex
	go reader.drain()
	return reader
}

func (r *bufReader) drain() {
	for {
		n, err := r.reader.Read(r.drainBuf)
		r.Lock()
		if err != nil {
			r.err = err
		} else {
			r.buf.Write(r.drainBuf[0:n])
		}
		r.wait.Signal()
		r.Unlock()
		if err != nil {
			break
		}
	}
}

func (r *bufReader) Read(p []byte) (n int, err error) {
	r.Lock()
	defer r.Unlock()
	for {
		n, err = r.buf.Read(p)
		if n > 0 {
			return n, err
		}
		if r.err != nil {
			return 0, r.err
		}
		r.wait.Wait()
	}
}

func (r *bufReader) Close() error {
	closer, ok := r.reader.(io.ReadCloser)
	if !ok {
		return nil
	}
	return closer.Close()
}
