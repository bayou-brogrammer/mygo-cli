package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// LogLevel represents the severity level of a log message
type LogLevel int

const (
	// LevelDebug for detailed troubleshooting information
	LevelDebug LogLevel = iota
	// LevelInfo for general operational information
	LevelInfo
	// LevelWarn for warning conditions
	LevelWarn
	// LevelError for error conditions
	LevelError
	// LevelFatal for critical errors that cause program termination
	LevelFatal
)

var (
	// Default logger instances for each level
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger

	// Current minimum log level
	currentLevel LogLevel = LevelInfo

	// Console output for debug and info logs
	consoleOutput io.Writer = os.Stdout

	// File output for warnings, errors, and fatal logs
	fileOutput io.Writer

	// Log file path
	logFilePath string

	// Log file handle
	logFile *os.File
)

// Init initializes the logger with the specified minimum level
func Init(level LogLevel) {
	currentLevel = level

	// Create log directory if it doesn't exist
	logDir := filepath.Join(os.Getenv("HOME"), ".config", "mygo", "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create log directory: %v\n", err)
	}

	// Set log file path
	logFilePath = filepath.Join(logDir, "mygo.log")

	// Open log file for warnings, errors, and fatal logs
	var err error
	logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open log file: %v\n", err)
		// Fallback to stderr for file logging if file can't be opened
		fileOutput = os.Stderr
	} else {
		fileOutput = logFile
	}

	setupLoggers()
}

// setupLoggers configures the logger instances
func setupLoggers() {
	// Debug and Info go to console
	debugLogger = log.New(consoleOutput, "DEBUG: ", log.Ldate|log.Ltime)
	infoLogger = log.New(consoleOutput, "INFO: ", log.Ldate|log.Ltime)

	// Warn, Error, and Fatal go to log file
	warnLogger = log.New(fileOutput, "WARN: ", log.Ldate|log.Ltime)
	errorLogger = log.New(fileOutput, "ERROR: ", log.Ldate|log.Ltime)
	fatalLogger = log.New(fileOutput, "FATAL: ", log.Ldate|log.Ltime)
}

// Close closes the log file
func Close() {
	if logFile != nil {
		logFile.Close()
		logFile = nil
	}
}

// ParseLevel converts a string level to LogLevel
func ParseLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn", "warning":
		return LevelWarn
	case "error":
		return LevelError
	case "fatal":
		return LevelFatal
	default:
		return LevelInfo
	}
}

// Debug logs a debug message if the current level allows it
func Debug(format string, v ...interface{}) {
	if currentLevel <= LevelDebug {
		debugLogger.Printf(format, v...)
	}
}

// Info logs an informational message if the current level allows it
func Info(format string, v ...interface{}) {
	if currentLevel <= LevelInfo {
		infoLogger.Printf(format, v...)
	}
}

// Warn logs a warning message if the current level allows it
func Warn(format string, v ...interface{}) {
	if currentLevel <= LevelWarn {
		warnLogger.Printf(format, v...)
	}
}

// Error logs an error message if the current level allows it
func Error(format string, v ...interface{}) {
	if currentLevel <= LevelError {
		errorLogger.Printf(format, v...)
	}
}

// Fatal logs a fatal message and exits the program
func Fatal(format string, v ...interface{}) {
	if currentLevel <= LevelFatal {
		fatalLogger.Printf(format, v...)
		os.Exit(1)
	}
}

// String returns the string representation of a LogLevel
func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	default:
		return fmt.Sprintf("LogLevel(%d)", l)
	}
}
