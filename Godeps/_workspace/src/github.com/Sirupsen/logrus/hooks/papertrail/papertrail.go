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

package logrus_papertrail

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
)

const (
	format = "Jan 2 15:04:05"
)

// PapertrailHook to send logs to a logging service compatible with the Papertrail API.
type PapertrailHook struct {
	Host    string
	Port    int
	AppName string
	UDPConn net.Conn
}

// NewPapertrailHook creates a hook to be added to an instance of logger.
func NewPapertrailHook(host string, port int, appName string) (*PapertrailHook, error) {
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", host, port))
	return &PapertrailHook{host, port, appName, conn}, err
}

// Fire is called when a log event is fired.
func (hook *PapertrailHook) Fire(entry *logrus.Entry) error {
	date := time.Now().Format(format)
	payload := fmt.Sprintf("<22> %s %s: [%s] %s", date, hook.AppName, entry.Level, entry.Message)

	bytesWritten, err := hook.UDPConn.Write([]byte(payload))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to send log line to Papertrail via UDP. Wrote %d bytes before error: %v", bytesWritten, err)
		return err
	}

	return nil
}

// Levels returns the available logging levels.
func (hook *PapertrailHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}
