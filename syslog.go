// Copyright 2021 xgfone
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build !windows,!plan9

package log

import (
	"bytes"
	"log/syslog"
)

type syslogWriter struct {
	w *syslog.Writer
}

func (s syslogWriter) WriteLevel(level Level, p []byte) (n int, err error) {
	v := string(bytes.TrimSpace(p))
	if level < LvlDebug {
		err = s.w.Debug(v)
	} else if level < LvlInfo {
		err = s.w.Info(v)
	} else if level < LvlWarn {
		err = s.w.Notice(v)
	} else if level < LvlError {
		err = s.w.Warning(v)
	} else if level < LvlFatal {
		err = s.w.Err(v)
	} else {
		err = s.w.Emerg(v)
	}

	if err == nil {
		n = len(p)
	}

	return
}

// Close implements io.Closer.
func (s syslogWriter) Close() error {
	return s.w.Close()
}

// SyslogWriter opens a connection to the system syslog daemon
// by calling syslog.New and writes all logs to it.
func SyslogWriter(priority syslog.Priority, tag string) (Writer, error) {
	w, err := syslog.New(priority, tag)
	if err != nil {
		return nil, err
	}
	return syslogWriter{w}, nil
}

// SyslogNetWriter opens a connection to a log daemon over the network
// and writes all logs to it.
func SyslogNetWriter(net, addr string, priority syslog.Priority, tag string) (Writer, error) {
	w, err := syslog.Dial(net, addr, priority, tag)
	if err != nil {
		return nil, err
	}
	return syslogWriter{w}, nil
}
