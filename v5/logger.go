package v5

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	logToConsole  bool
	logToFile     bool
	logDir        string
	fileLogger    *log.Logger
	consoleLogger *log.Logger
	logFile       *os.File
	lastLogDate   string
	logLock       sync.Mutex
)

// InitLogger sets up the logging system.
// Enable or disable console/file logging independently.
// logDir specifies the directory for daily log files.
func InitLogger(enableConsole bool, enableFile bool, dir string) error {
	logToConsole = enableConsole
	logToFile = enableFile
	logDir = dir
	lastLogDate = ""

	if logToFile {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return fmt.Errorf("failed to create log directory: %w", err)
		}
	}
	consoleLogger = log.New(os.Stdout, "", 0)
	return nil
}

// LogInfof logs an informational message.
func LogInfof(format string, args ...interface{}) {
	logf("INFO", format, args...)
}

// LogWarnf logs a warning message.
func LogWarnf(format string, args ...interface{}) {
	logf("WARN", format, args...)
}

// LogErrorf logs an error message.
func LogErrorf(format string, args ...interface{}) {
	logf("ERROR", format, args...)
}

// CloseLogger flushes and closes the file logger if open.
func CloseLogger() {
	logLock.Lock()
	defer logLock.Unlock()

	if logFile != nil {
		logFile.Close()
		logFile = nil
	}
}

// logf is the internal shared formatter and writer for all log levels.
func logf(level string, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fullMessage := fmt.Sprintf("[%s] [TSL] [%s] %s", timestamp, level, message)

	logLock.Lock()
	defer logLock.Unlock()

	// Console color output
	if logToConsole {
		colored := colorWrap(level, fullMessage)
		consoleLogger.Println(colored)
	}

	// File logging
	if logToFile {
		rotateLogFileIfNeeded()
		if fileLogger != nil {
			fileLogger.Println(fullMessage)
		}
	}
}

// rotateLogFileIfNeeded checks if the current log file matches today and creates a new one if necessary.
func rotateLogFileIfNeeded() {
	currentDate := time.Now().Format("2006-01-02")
	if currentDate == lastLogDate && logFile != nil {
		return
	}

	if logFile != nil {
		logFile.Close()
	}

	filename := fmt.Sprintf("TSLLog-%s.log", currentDate)
	fullPath := filepath.Join(logDir, filename)

	f, err := os.OpenFile(fullPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		consoleLogger.Printf("[TSL] [ERROR] Failed to open log file: %v", err)
		logToFile = false
		return
	}

	logFile = f
	fileLogger = log.New(logFile, "", 0)
	lastLogDate = currentDate
}

// colorWrap applies ANSI color codes for console output based on log level.
func colorWrap(level string, msg string) string {
	switch level {
	case "ERROR":
		return "\033[31m" + msg + "\033[0m"
	case "WARN":
		return "\033[33m" + msg + "\033[0m"
	default:
		return msg
	}
}
