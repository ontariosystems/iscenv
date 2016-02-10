/*
Copyright 2016 Ontario Systems

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

package iscenv

import (
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/creack/termios/raw"
	"github.com/creack/termios/win"
)

type TTYSizeChangeHandler func(height, width int)

func NewRawTTYStdin() (tty *rawTTYStdin, err error) {
	rts := &rawTTYStdin{File: os.Stdin}
	if rts.IsTerminal() {
		rts.origState, err = raw.MakeRaw(rts.Fd())
		if err != nil {
			return nil, err
		}

	}

	return rts, nil
}

type rawTTYStdin struct {
	origState *raw.Termios
	sigchan   chan os.Signal
	*os.File
}

func (rts *rawTTYStdin) Fd() uintptr {
	return rts.File.Fd()
}

func (rts *rawTTYStdin) IsTerminal() bool {
	return terminal.IsTerminal(int(rts.Fd()))
}

func (rts *rawTTYStdin) Close() error {
	if rts.sigchan != nil {
		signal.Stop(rts.sigchan)
	}

	if rts.origState != nil {
		raw.TcSetAttr(rts.Fd(), rts.origState)
	}

	return rts.File.Close()
}

func (rts *rawTTYStdin) MonitorTTYSize(handler TTYSizeChangeHandler) {
	// Call the handler for the initial TTY size
	rts.resizeTTY(handler)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGWINCH)
	go func() {
		for range sigchan {
			rts.resizeTTY(handler)
		}
	}()
}

func (rts *rawTTYStdin) resizeTTY(handler TTYSizeChangeHandler) {
	height, width := rts.GetTTYSize()
	// Could not determine the height & width so, don't call the handler
	if height == 0 && width == 0 {
		return
	}

	handler(height, width)
}

func (rts *rawTTYStdin) GetTTYSize() (height, width int) {
	if rts.IsTerminal() {
		ws, err := win.GetWinsize(rts.Fd())
		if err == nil {
			return int(ws.Height), int(ws.Width)
		}
	}

	return 0, 0
}
