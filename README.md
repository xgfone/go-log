# Go Log [![Build Status](https://github.com/xgfone/go-log/actions/workflows/go.yml/badge.svg)](https://github.com/xgfone/go-log/actions/workflows/go.yml) [![GoDoc](https://pkg.go.dev/badge/github.com/xgfone/go-log)](https://pkg.go.dev/github.com/xgfone/go-log) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=flat-square)](https://raw.githubusercontent.com/xgfone/go-log/master/LICENSE)

Provide a simple, flexible, extensible, powerful and structured logger based on the level, which has done the better balance between the flexibility and the performance. It is inspired by [log15](https://github.com/inconshreveable/log15), [logrus](https://github.com/sirupsen/logrus), [go-kit](https://github.com/go-kit/kit) and [zerolog](github.com/rs/zerolog), which collects the log message with the key-value contexts, encodes them into the buffer, then writes the encoded log from the buffer into the underlying writer.


## Features

- Support `Go1.7+`.
- Compatible with the stdlib `log.Printf`.
- The better performance, see [Benchmark](#performance).
    - Lazy evaluation of expensive operations.
    - Avoid to allocate the memory on heap as far as possible.
    - Encode in real time or pre-encode the key-value contexts into the buffer cache.
- Simple, Flexible, Extensible, Powerful and Structured.
- Support to customize the log encoder and writer.
- Provide the simple and easy-used api interface.


## Example

```go
package main

import (
    "errors"
    "flag"

    "github.com/xgfone/go-log"
)

var logfile string
var loglevel string

func logError(err error, msg string, kvs ...interface{}) {
    if err == nil {
        return
    }
    log.Level(log.LvlError, 1).Kvs(kvs...).Kv("err", err).Printf(msg)
}

func main() {
    // Parse the CLI options.
    flag.StringVar(&logfile, "logfile", "", "The log file path, default to stderr.")
    flag.StringVar(&loglevel, "loglevel", "info", "The log level, such as debug, info, etc.")
    flag.Parse()

    // Configure the logger.
    writer := log.FileWriter(logfile, "100M", 100)
    defer writer.Close()
    log.SetWriter(writer)
    log.SetLevel(log.ParseLevel(loglevel))

    // Emit the log.
    log.Debug().Kv("key1", "value1").Print("msg1") // no log output.
    log.Info().Kv("key2", "value2").Print("msg2")
    log.Level(log.LvlInfo, 0).Kv("key3", "value3").Printf("msg3")
    logError(nil, "msg4", "key4", "value4", "key5", 555, "key6", "value6")
    logError(errors.New("error"), "msg7", "key8", 888, "key9", "value9")

    // For Clild Logger
    child1Logger := log.WithName("child1")
    child2Logger := child1Logger.WithName("child2")
    child1Logger.Info().Kv("ckey1", "cvalue1").Print("msg8")
    child2Logger.Info().Kv("ckey2", "cvalue2").Printf("msg9")

    // $ go run main.go
    // {"t":"2021-12-17T00:04:44.8609884+08:00","lvl":"info","caller":"main.go:34:main","key2":"value2","msg":"msg2"}
    // {"t":"2021-12-17T00:04:44.8660577+08:00","lvl":"info","caller":"main.go:35:main","key3":"value3","msg":"msg3"}
    // {"t":"2021-12-17T00:04:44.8671207+08:00","lvl":"error","caller":"main.go:37:main","key8":888,"key9":"value9","err":"error","msg":"msg7"}
    // {"t":"2021-12-17T00:04:44.8671207+08:00","lvl":"info","logger":"child1","caller":"main.go:42:main","ckey1":"cvalue1","msg":"msg8"}
    // {"t":"2021-12-17T00:04:44.8678731+08:00","lvl":"info","logger":"child1.child2","caller":"main.go:43:main","ckey2":"cvalue2","msg":"msg9"}
}
```


### `logr`

```go
// logr.go
package main

import (
    "fmt"

    "github.com/go-logr/logr"
    "github.com/xgfone/go-log"
)

// NewLogSink returns a logr sink based on the key-value logger.
func NewLogSink(logger log.Logger) logr.LogSink {
    return &logSink{logger: logger}
}

const maxLevel = log.LvlWarn - log.LvlInfo - 1

type logSink struct {
    logger log.Logger
}

func (l *logSink) Init(info logr.RuntimeInfo) {
    l.logger = l.logger.WithDepth(info.CallDepth + 1)
}

func (l *logSink) Enabled(level int) bool {
    if level > maxLevel {
        panic(fmt.Errorf("invalid level '%d': only allow [0, %d]", level, maxLevel))
    }
    return l.logger.Enabled(log.LvlInfo + level)
}

func (l *logSink) Info(level int, msg string, keysAndValues ...interface{}) {
    if level > maxLevel {
        panic(fmt.Errorf("invalid level '%d': only allow [0, %d]", level, maxLevel))
    }
    l.logger.Level(log.LvlInfo+level, l.logger.Depth()+1).Kvs(keysAndValues...).Printf(msg)
}

func (l *logSink) Error(err error, msg string, keysAndValues ...interface{}) {
    l.logger.Error().Kvs(keysAndValues...).Kv("err", err).Printf(msg)
}

func (l *logSink) WithName(name string) logr.LogSink {
    return &logSink{l.logger.WithName(name)}
}

func (l *logSink) WithValues(keysAndValues ...interface{}) logr.LogSink {
    return &logSink{l.logger.WithContexts(keysAndValues...)}
}

func (l *logSink) WithCallDepth(depth int) logr.LogSink {
    return &logSink{l.logger.WithDepth(depth + 2)}
}
```

```go
// main.go
package main

import (
    "errors"

    "github.com/go-logr/logr"
    "github.com/xgfone/go-log"
)

func logIfErr(logger logr.Logger, err error, msg string, kvs ...interface{}) {
    if err != nil {
        logger.Error(err, msg, kvs...)
    }
}

func main() {
    _logger := log.New("test").
        WithHooks(log.Caller("caller")). // Add the caller context
        WithLevel(log.LvlInfo + 3)       // Only output the logs that V is not less than 3

    logger := logr.New(NewLogSink(_logger))
    logger.Info("msg1", "k11", "v11", "k12", "v12") // The log is not be output.
    logger.Error(errors.New("error"), "msg2", "k2", "v2")

    logger = logger.V(6) // V must be between 0 and 19, that's, [0, 19].
    logger.Info("msg3", "k3", "v3")
    logger.Error(errors.New("error"), "msg4", "k4", "v4")

    logger = logger.WithName("name")
    logger.Info("msg5", "k5", "v5")
    logger.Error(errors.New("error"), "msg6", "k6", "v6")

    logger = logger.WithValues("k0", "v0")
    logger.Info("msg7", "k7", "v7")
    logger.Error(errors.New("error"), "msg8", "k8", "v8")

    logger = logger.WithCallDepth(1)
    logIfErr(logger, errors.New("error"), "msg9", "k9", "v9")
    logIfErr(logger, nil, "msg10", "k10", "v10")

    // $ go run logr.go main.go
    // {"t":"2021-12-17T00:16:10.1478129+08:00","lvl":"error","logger":"test","caller":"main.go:23:main","k2":"v2","err":"error","msg":"msg2"}
    // {"t":"2021-12-17T00:16:10.1535681+08:00","lvl":"info6","logger":"test","k3":"v3","msg":"msg3"}
    // {"t":"2021-12-17T00:16:10.1541601+08:00","lvl":"error","logger":"test","caller":"main.go:27:main","k4":"v4","err":"error","msg":"msg4"}
    // {"t":"2021-12-17T00:16:10.1546859+08:00","lvl":"info6","logger":"test.name","k5":"v5","msg":"msg5"}
    // {"t":"2021-12-17T00:16:10.1546859+08:00","lvl":"error","logger":"test.name","caller":"main.go:31:main","k6":"v6","err":"error","msg":"msg6"}
    // {"t":"2021-12-17T00:16:10.1552482+08:00","lvl":"info6","logger":"test.name","k0":"v0","k7":"v7","msg":"msg7"}
    // {"t":"2021-12-17T00:16:10.1552482+08:00","lvl":"error","logger":"test.name","k0":"v0","caller":"main.go:35:main","k8":"v8","err":"error","msg":"msg8"}
    // {"t":"2021-12-17T00:16:10.1558789+08:00","lvl":"error","logger":"test.name","k0":"v0","caller":"main.go:38:main","k9":"v9","err":"error","msg":"msg9"}
}
```


### Encoder

```go
type Encoder interface {
    // Start starts to encode the log record into the buffer dst.
    Start(dst []byte, loggerName string, level int) []byte

    // Encode encodes the key-value with the stack depth into the buffer dst.
    Encode(dst []byte, key string, value interface{}) []byte

    // End ends to encode the log record with the message into the buffer dst.
    End(dst []byte, msg string) []byte
}
```

This pakcage has implemented the JSON encoder `JSONEncoder`, but you can customize yourself, such as `TextEncoder`.


### Writer

The logger uses the stdlib `io.Writer` interface as the log writer.

In order to support to write the leveled log, you can provide a `LevelWriter` to the log engine, which prefers to try to use `LevelWriter` to write the log into it.
```go
type LevelWriter interface {
    WriteLevel(level int, data []byte) (n int, err error)
    io.Writer
}
```

The package provides an additional writer based on the file, that's, `FileWriter`.


### Sampler

The logger engine provides the sampler policy for each logger to filter the log message by the logger name and level during the program is running.
```go
type Sampler interface {
    // Sample reports whether the log message should be sampled.
    // If the log message should be sampled, return true. Or, return false,
    // that's, the log message will be discarded.
    Sample(loggerName string, level int) bool
}
```

Notice: in order to switch the level of all the loggers once, you maybe use the global level function `SetGlobalLevel`, such as `SetGlobalLevel(LvlError)`, which will disable all the log messages whose level is lower than `LvlError`.


### Lazy evaluation
The logger provides the hook `Hook` to support the Lazy evaluation.
```go
type Hook interface {
    Run(logger Logger, loggerName string, level int, depth int)
}
```

The package provides a dynamic key-value context `Caller` to calculate the file and line where the caller is.
```go
package main

import "github.com/xgfone/go-log"

func main() {
    logger := log.New("root").WithHooks(log.Caller("caller"))
    logger.Info().Kv("key", "value").Printf("msg")

    // $ go run main.go
    // {"t":"2021-12-12T15:09:41.6890462+08:00","lvl":"info","logger":"root","caller":"main.go:7:main","key":"value","msg":"msg"}
}
```

Not only the lazy evaluation, but the hook is also used to do others, such as the counter of the level logs.


## Performance

The log framework itself has no any performance costs and the key of the bottleneck is the encoder.

```
HP Laptop 14s-dr2014TU
go: 1.17.3
goos: windows
goarch: amd64
cpu: 11th Gen Intel(R) Core(TM) i7-1165G7 @ 2.80GHz
```

**Benchmark Package:**
|                Function                |      ops      | ns/op | bytes/opt | allocs/op
|----------------------------------------|--------------:|------:|-----------|----------
|BenchmarkJSONEncoderDisabled-8          | 325, 556, 422 | 3.649 |     0     |    0
|BenchmarkJSONEncoderEmpty-8             |  71, 245, 855 | 17.71 |     0     |    0
|BenchmarkJSONEncoderInfo-8              |  64, 453, 407 | 17.84 |     0     |    0
|BenchmarkJSONEncoderWith8Contexts-8     |  63, 589, 971 | 17.87 |     0     |    0
|BenchmarkJSONEncoderWith8KeyValues-8    |   9, 351, 409 | 121.9 |    128    |    8
|BenchmarkJSONEncoderWithOptimized8KVs-8 |  14, 620, 470 | 78.71 |     2     |    1
