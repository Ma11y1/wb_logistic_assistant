package logger

import (
	"errors"
	"fmt"
	"os"
	"time"
)

var timeFormat = "2006.01.02 15:04:05 -0700"
var fatalHandler = func() {}

func SetTimeFormat(format string) error {
	if format == "" {
		return errors.New("time format is empty")
	}
	mtx.Lock()
	timeFormat = format
	mtx.Unlock()
	return nil
}

func SetFatalHandler(handler func()) {
	if handler != nil {
		fatalHandler = handler
	}
}

func Log(level Level, location, msg string) {
	mtx.RLock()
	defer mtx.RUnlock()
	outs := getOutputs(level)
	ts := time.Now().Format(timeFormat)
	for i := 0; i < len(outs); i++ {
		_, _ = fmt.Fprintf(outs[i], "%s: [%s] %s: %s\n", ts, level, location, msg)
	}
	if level == FATAL {
		fatalHandler()
		os.Exit(1)
	}
}

func Logf(level Level, location, format string, args ...interface{}) {
	Log(level, location, fmt.Sprintf(format, args...))
}
