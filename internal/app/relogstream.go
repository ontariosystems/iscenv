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

package app

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	log "github.com/Sirupsen/logrus"
)

// RelogStream will given a JSON log stream, relog all of the message into the current process
// extraFields contains additional fields you wish to add to the log output.  These will overwrite any fields from the original log message.
// preserveTimestamps will use the timestamps from the original log messages instead of the current time
func RelogStream(extraFields log.Fields, preserveTimestamp bool, r io.Reader) {
	decoder := json.NewDecoder(r)
	for {
		var rlog map[string]interface{}
		if err := decoder.Decode(&rlog); err == nil {
			Relog(extraFields, preserveTimestamp, rlog)
		} else if err == io.EOF {
			return
		} else {
			ErrorLogger(log.WithFields(extraFields), err).Warn("Failed to parse streamed log message")
		}
	}
}

// Relog will relog a single log message that has been Unmarshaled into a map[string]interface{}
// extraFields contains additional fields you wish to add to the log output.  These will overwrite any fields from the original log message.
// preserveTimestamps will use the timestamps from the original log messages instead of the current time
func Relog(extraFields log.Fields, preserveTimestamp bool, rlog map[string]interface{}) {
	l := log.NewEntry(log.StandardLogger())
	for key, value := range rlog {
		switch key {
		case "time":
			if preserveTimestamp {
				if ts, err := parseTime(value); err == nil {
					l = l.WithField("overrideTime", ts)
				} else {
					log.WithError(err).Warn("Could not parse time, not preserving timestamp")
				}
			}
		case "level", "msg":
			// Skip
		default:
			l = l.WithField(key, value)
		}
	}

	var level log.Level
	var err error
	if levelStr, ok := rlog["level"].(string); ok {
		if level, err = log.ParseLevel(levelStr); err != nil {
			level = log.InfoLevel
			l.WithField("origLevel", levelStr).Warn("Unknown log level, using info")
		}
	} else {
		level = log.InfoLevel
		l.WithField("origLevel", rlog["level"]).Warn("Remote log level was not a string, using info")
	}

	msg, ok := rlog["msg"].(string)
	if !ok {
		log.WithField("msg", msg).Error("Remote log message was not a string, skipping")
		return
	}

	// Ensure that these fields overwrite any from the message
	l.WithFields(extraFields)

	switch level {
	case log.DebugLevel:
		l.Debug(msg)
	case log.InfoLevel:
		l.Info(msg)
	case log.WarnLevel:
		l.Warn(msg)
	case log.ErrorLevel:
		l.Error(msg)
	default:
		// At this point it means we are Fatal or Panic, we don't really want to log either of those levels as their are side effects
		l.WithField("origLevel", level.String()).Error(msg)
	}
}

func parseTime(ti interface{}) (time.Time, error) {
	switch t := ti.(type) {
	case string:
		return time.Parse(time.RFC3339, t)
	default:
		return time.Time{}, fmt.Errorf("Remote time field was not a string")
	}
}
