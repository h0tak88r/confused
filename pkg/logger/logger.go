package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
)

// LogLevel represents the logging level
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// Logger represents a structured logger
type Logger struct {
	level    LogLevel
	verbose  bool
	file     *os.File
	filePath string
}

// New creates a new logger instance
func New(level LogLevel, verbose bool, logFile string) (*Logger, error) {
	logger := &Logger{
		level:   level,
		verbose: verbose,
	}

	if logFile != "" {
		// Ensure directory exists
		dir := filepath.Dir(logFile)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}

		logger.file = file
		logger.filePath = logFile
	}

	return logger, nil
}

// Close closes the logger file
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// log writes a log message with the given level
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	levelStr := l.getLevelString(level)
	message := fmt.Sprintf(format, args...)
	
	logMessage := fmt.Sprintf("[%s] [%s] %s", timestamp, levelStr, message)

	// Write to console with colors
	l.writeToConsole(level, logMessage)

	// Write to file if available
	if l.file != nil {
		fmt.Fprintln(l.file, logMessage)
	}
}

// writeToConsole writes the message to console with appropriate colors
func (l *Logger) writeToConsole(level LogLevel, message string) {
	switch level {
	case DEBUG:
		if l.verbose {
			color.New(color.FgCyan).Println(message)
		}
	case INFO:
		color.New(color.FgGreen).Println(message)
	case WARN:
		color.New(color.FgYellow).Println(message)
	case ERROR:
		color.New(color.FgRed).Println(message)
	case FATAL:
		color.New(color.FgRed, color.Bold).Println(message)
	default:
		fmt.Println(message)
	}
}

// getLevelString returns the string representation of the log level
func (l *Logger) getLevelString(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
	os.Exit(1)
}

// SetLevel sets the logging level
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// SetVerbose sets the verbose mode
func (l *Logger) SetVerbose(verbose bool) {
	l.verbose = verbose
}

// GetFileWriter returns a writer for the log file
func (l *Logger) GetFileWriter() io.Writer {
	if l.file != nil {
		return l.file
	}
	return os.Stdout
}
