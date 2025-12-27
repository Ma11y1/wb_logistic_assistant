package logger

import (
	"errors"
	"io"
	"sync"
)

var (
	outputs      = make([]io.Writer, 0)
	errorOutputs = make([]io.Writer, 0)
	mtx          = sync.RWMutex{}
)

func AddOutput(writer io.Writer) error {
	if writer == nil {
		return errors.New("writer is nil")
	}
	mtx.Lock()
	defer mtx.Unlock()
	for _, w := range outputs {
		if w == writer {
			return errors.New("writer already exists")
		}
	}
	outputs = append(outputs, writer)
	return nil
}

func RemoveOutput(writer io.Writer) {
	if writer == nil {
		return
	}
	mtx.Lock()
	defer mtx.Unlock()
	for i, w := range outputs {
		if w == writer {
			outputs = append(outputs[:i], outputs[i+1:]...)
			return
		}
	}
}

func AddOutputErr(writer io.Writer) error {
	if writer == nil {
		return errors.New("writer is nil")
	}
	mtx.Lock()
	defer mtx.Unlock()
	for _, o := range errorOutputs {
		if o == writer {
			return errors.New("writer already exists")
		}
	}
	errorOutputs = append(errorOutputs, writer)
	return nil
}

func RemoveOutputErr(writer io.Writer) {
	if writer == nil {
		return
	}
	mtx.Lock()
	defer mtx.Unlock()
	for i, o := range errorOutputs {
		if o == writer {
			errorOutputs = append(errorOutputs[:i], errorOutputs[i+1:]...)
			return
		}
	}
}

func getOutputs(level Level) []io.Writer {
	switch level {
	case FATAL, ERROR, WARN:
		return errorOutputs
	default:
		return outputs
	}
}
