package v5

import (
	"bytes"
	"io"
	"log"
	"os"
	"sync"
)

var (
	logger       *log.Logger
	logFile      *os.File
	memoryBuffer *circularBuffer
)

func InitLogger(toFile bool, filepath string) error {
	var output io.Writer

	memoryBuffer = newCircularBuffer(50)

	if toFile {
		f, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		logFile = f
		output = io.MultiWriter(os.Stdout, f, memoryBuffer)
	} else {
		output = io.MultiWriter(os.Stdout, memoryBuffer)
	}

	logger = log.New(output, "", log.LstdFlags|log.Lshortfile)
	return nil
}

func Logf(format string, args ...interface{}) {
	if logger != nil {
		logger.Printf(format, args...)
	}
}

func CloseLogger() {
	if logFile != nil {
		logFile.Close()
	}
}

func SnapshotLog() []string {
	if memoryBuffer != nil {
		return memoryBuffer.Snapshot()
	}
	return nil
}

type circularBuffer struct {
	mu     sync.Mutex
	lines  []string
	max    int
	cursor int
	full   bool
}

func newCircularBuffer(size int) *circularBuffer {
	return &circularBuffer{
		lines: make([]string, size),
		max:   size,
	}
}

func (cb *circularBuffer) Write(p []byte) (n int, err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.lines[cb.cursor] = string(bytes.TrimSpace(p))
	cb.cursor = (cb.cursor + 1) % cb.max
	if cb.cursor == 0 {
		cb.full = true
	}

	return len(p), nil
}

func (cb *circularBuffer) Snapshot() []string {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	var out []string
	if cb.full {
		out = append(out, cb.lines[cb.cursor:]...)
		out = append(out, cb.lines[:cb.cursor]...)
	} else {
		out = append(out, cb.lines[:cb.cursor]...)
	}
	return out
}
