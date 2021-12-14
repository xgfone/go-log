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
    ```go
    type Logger interface {
        // Enabled reports whether the logger is enabled.
        Enabled() bool

        // Kv and Kvs append the key-value contexts and return the logger itself.
        Kv(key string, value interface{}) Logger
        Kvs(kvs ...interface{}) Logger

        // Print and Printf log the message and end the logger.
        Printf(msg string, args ...interface{})
        Print(args ...interface{})
    }
    ```


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
    log.Log(log.LvlError, 1).Kvs(kvs...).Kv("err", err).Printf(msg)
}

func main() {
    // Parse the CLI options.
    flag.StringVar(&logfile, "logfile", "", "The log file path, default to stderr.")
    flag.StringVar(&loglevel, "loglevel", "info", "The log level, such as debug, info, etc.")
    flag.Parse()

    // Configure the logger.
    writer := log.FileWriter(logfile, "100M", 100)
    log.SetWriter(writer).SetLevel(log.ParseLevel(loglevel))
    defer writer.Close()

    // Emit the log.
    log.Print("msg1")
    log.Printf("msg%d", 2)
    log.Kv("key1", "value1").Print("msg3")
    log.Debug().Kv("key2", "value2").Print("msg4") // no log output.
    log.Info().Kv("key3", "value3").Print("msg5")
    log.Log(log.LvlInfo, 0).Kv("key4", "value4").Printf("msg6")
    logError(nil, "msg7", "key5", "value5", "key6", 666, "key7", "value7")
    logError(errors.New("error"), "msg8", "key8", 888, "key9", "value9")

    // For Clild Logger
    child1Logger := log.WithName("child1")
    child2Logger := child1Logger.New("child2")
    child1Logger.Kv("ckey1", "cvalue1").Print("msg9")
    child2Logger.Printf("msg10")

    // $ go run main.go
    // {"t":"2021-12-12T11:41:11.2844234+08:00","lvl":"info","caller":"main.go:32:main","msg":"msg1"}
    // {"t":"2021-12-12T11:41:11.2918549+08:00","lvl":"info","caller":"main.go:33:main","msg":"msg2"}
    // {"t":"2021-12-12T11:41:11.2918549+08:00","lvl":"info","caller":"main.go:34:main","key1":"value1","msg":"msg3"}
    // {"t":"2021-12-12T11:41:11.2918549+08:00","lvl":"info","caller":"main.go:36:main","key3":"value3","msg":"msg5"}
    // {"t":"2021-12-12T11:41:11.2918549+08:00","lvl":"info","caller":"main.go:37:main","key4":"value4","msg":"msg6"}
    // {"t":"2021-12-12T11:41:11.2918549+08:00","lvl":"error","caller":"main.go:39:main","key8":888,"key9":"value9","err":"error","msg":"msg8"}
    // {"t":"2021-12-12T12:22:15.2466635+08:00","lvl":"info","logger":"child1","caller":"main.go:44:main","ckey1":"cvalue1","msg":"msg9"}
    // {"t":"2021-12-12T12:22:15.2466635+08:00","lvl":"info","logger":"child1.child2","caller":"main.go:45:main","msg":"msg10"}
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
func NewLogSink(logger *log.Engine) logr.LogSink {
    return &logSink{logger: logger}
}

const maxLevel = log.LvlWarn - log.LvlInfo - 1

type logSink struct {
    logger *log.Engine
}

func (l logSink) Init(info logr.RuntimeInfo) {
    l.logger.SetDepth(info.CallDepth + 1)
}

func (l logSink) Enabled(level int) bool {
    if level > maxLevel {
        panic(fmt.Errorf("invalid level '%d': only allow [0, %d]", level, maxLevel))
    }
    return l.logger.Enable(log.LvlInfo + level)
}

func (l logSink) Info(level int, msg string, keysAndValues ...interface{}) {
    if level > maxLevel {
        panic(fmt.Errorf("invalid level '%d': only allow [0, %d]", level, maxLevel))
    }
    l.logger.Level(log.LvlInfo + level).Kvs(keysAndValues...).Printf(msg)
}

func (l logSink) Error(err error, msg string, keysAndValues ...interface{}) {
    l.logger.Error().Kvs(keysAndValues...).Kv("err", err).Printf(msg)
}

func (l logSink) WithName(name string) logr.LogSink {
    return logSink{l.logger.New(name)}
}

func (l logSink) WithValues(keysAndValues ...interface{}) logr.LogSink {
    return logSink{l.logger.Clone().AppendCtxs(keysAndValues...)}
}

func (l logSink) WithCallDepth(depth int) logr.LogSink {
    return logSink{l.logger.Clone().SetDepth(depth + 2)}
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
    engine := log.New("test").AddHooks(log.Caller("caller"))
    engine.SetLevel(log.LvlInfo + 3) // Only output the logs that V is not less than 3

    logger := logr.New(NewLogSink(engine))
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
    // {"t":"2021-12-14T23:31:04.457377+08:00","lvl":"error","logger":"test","caller":"main.go:22:main","k2":"v2","err":"error","msg":"msg2"}
    // {"t":"2021-12-14T23:31:04.4637531+08:00","lvl":"info6","logger":"test","caller":"main.go:25:main","k3":"v3","msg":"msg3"}
    // {"t":"2021-12-14T23:31:04.4637531+08:00","lvl":"error","logger":"test","caller":"main.go:26:main","k4":"v4","err":"error","msg":"msg4"}
    // {"t":"2021-12-14T23:31:04.4643901+08:00","lvl":"info6","logger":"test.name","caller":"main.go:29:main","k5":"v5","msg":"msg5"}
    // {"t":"2021-12-14T23:31:04.4644614+08:00","lvl":"error","logger":"test.name","caller":"main.go:30:main","k6":"v6","err":"error","msg":"msg6"}
    // {"t":"2021-12-14T23:31:04.4648517+08:00","lvl":"info6","logger":"test.name","k0":"v0","caller":"main.go:33:main","k7":"v7","msg":"msg7"}
    // {"t":"2021-12-14T23:31:04.4648517+08:00","lvl":"error","logger":"test.name","k0":"v0","caller":"main.go:34:main","k8":"v8","err":"error","msg":"msg8"}
    // {"t":"2021-12-14T23:31:04.46535+08:00","lvl":"error","logger":"test.name","k0":"v0","caller":"main.go:37:main","k9":"v9","err":"error","msg":"msg9"}
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
    logger := log.New("root").AddHooks(log.Caller("caller"))
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
|                     Function                     |        ops       | ns/op  | bytes/opt | allocs/op
|--------------------------------------------------|-----------------:|-------:|-----------|----------
|BenchmarkLevelDisabled-8                          | 1, 000, 000, 000 |   0.76 |     0     |    0
|BenchmarkNothingEncoder-8                         |    150, 158, 604 |  11.00 |     0     |    0
|BenchmarkJSONEncoderWithoutContextsAndKeyValues-8 |     82, 356, 866 |  15.01 |     0     |    0
|BenchmarkJSONEncoderWith8Contexts-8               |     84, 787, 077 |  14.40 |     0     |    0
|BenchmarkJSONEncoderWith8KeyValues-8              |     26, 949, 334 |  45.83 |     0     |    0
