# go-log [![Build Status](https://travis-ci.org/xgfone/go-log.svg?branch=master)](https://travis-ci.org/xgfone/go-log) [![GoDoc](https://godoc.org/github.com/xgfone/go-log?status.svg)](http://godoc.org/github.com/xgfone/go-log) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=flat-square)](https://raw.githubusercontent.com/xgfone/go-log/master/LICENSE)

Provide a simple, flexible, extensible, powerful and structured logging tool based on the level, which has done the better balance between the flexibility and the performance. It is inspired by [log15](https://github.com/inconshreveable/log15), [logrus](https://github.com/sirupsen/logrus), [go-kit](https://github.com/go-kit/kit) and [zerolog](github.com/rs/zerolog).


## Features

- The better performance.
- Lazy evaluation of expensive operations.
- Simple, Flexible, Extensible, Powerful and Structured.
- Avoid to allocate the memory on heap as far as possible.
- Child loggers which inherit and add their own private context.
- Built-in support for logging to files, syslog, etc. See `Writer`.

The logger supports two kinds of the logger interfaces:
```go
// For Key-Value fields
Trace(msg string, fields ...Field)
Debug(msg string, fields ...Field)
Info(msg string, fields ...Field)
Warn(msg string, fields ...Field)
Error(msg string, fields ...Field)
Fatal(mst string, fields ...Field)

// For format string
Tracef(msgfmt string, args ...interface{})
Debugf(msgfmt string, args ...interface{})
Infof(msgfmt string, args ...interface{})
Warnf(msgfmt string, args ...interface{})
Errorf(msgfmt string, args ...interface{})
Fatalf(msgfmt string, args ...interface{})

// For Key-Value sequences
Traces(msg string, keyAndValues ...interface{})
Debugs(msg string, keyAndValues ...interface{})
Infos(msg string, keyAndValues ...interface{})
Warns(msg string, keyAndValues ...interface{})
Errors(msg string, keyAndValues ...interface{})
Fatals(mst string, keyAndValues ...interface{})
```


## Example

```go
package main

import "github.com/xgfone/go-log"

func main() {
	logger := log.New("name").WithLevel(log.LvlWarn)

	logger.Info("log msg", log.F("key1", "value1"), log.F("key2", "value2"))
	logger.Error("log msg", log.F("key1", "value1"), log.F("key2", "value2"))

	// Output:
	// {"t":"2021-05-28T22:00:00.092641+08:00","lvl":"ERROR","logger":"name","stack":"[main.go:9]","key1":"value1","key2":"value2","msg":"log msg"}
}
```

```go
package main

import "github.com/xgfone/go-log"

func main() {
	log.DefalutLogger.Level = log.LvlWarn

	// Emit the log with the fields.
	log.Info("log msg", log.F("key1", "value1"), log.F("key2", "value2"))
	log.Error("log msg", log.F("key1", "value1"), log.F("key2", "value2"))

	// Emit the log with key-values
	log.Infos("log msg", "key1", "value1", "key2", "value2")
	log.Errors("log msg", "key1", "value1", "key2", "value2")

	// Emit the log with the formatter.
	log.Infof("log %s", "msg")
	log.Errorf("log %s", "msg")

	// Output:
	// {"t":"2021-05-28T22:07:07.394835+08:00","lvl":"ERROR","stack":"[main.go:10]","key1":"value1","key2":"value2","msg":"log msg"}
	// {"t":"2021-05-28T22:07:07.395066+08:00","lvl":"ERROR","stack":"[main.go:14]","key1":"value1","key2":"value2","msg":"log msg"}
	// {"t":"2021-05-28T22:07:07.3951+08:00","lvl":"ERROR","stack":"[main.go:18]","msg":"log msg"}
}
```


### Encoder

```go
type Encoder interface {
	// Writer returns the writer.
	Writer() Writer

	// SetWriter resets the writer.
	SetWriter(Writer)

	// Encode encodes the log record and writes it into the writer.
	Encode(Record)
}
```

This pakcage has implemented the Nothing and JSON encoder, such as `NothingEncoder` and `JSONEncoder`.


### Writer

```go
type Writer interface {
	WriteLevel(level Level, data []byte) (n int, err error)
	io.Closer
}
```

There are some built-in writers, such as `DiscardWriter`, `LevelWriter`, `SafeWriter`, `SplitWriter`, `StreamWriter` and `FileWriter`.


### Lazy evaluation
`Field` supports the lazy evaluation, such as `F("key", func() interface{} { return "value" })`.


## Performance

The log framework itself has no any performance costs and the key of the bottleneck is the encoder.

```
Dell Vostro 3470
Intel Core i5-7400 3.0GHz
8GB DDR4 2666MHz
Windows 10
Go 1.16.4
```

**Benchmark Package:**

|               Function               |      ops      | ns/op | bytes/opt | allocs/op
|--------------------------------------|--------------:|------:|-----------|----------
|BenchmarkNothingEncoder-4             | 261, 674, 554 |   4.5 |     0     |    0
|BenchmarkJSONEncoderWithoutCtxField-4 |  11, 538, 594 |  98.0 |     0     |    0
|BenchmarkJSONEncoderWith10CtxFields-4 |   4, 109, 601 | 290.6 |     0     |    0
