package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	logDir      = "logs"
	logFile     = "logs/server.log"
	maxLogSize  = 10 // MB
	maxBackups  = 10 // Keep last 10 files
	maxAge      = 30 // days (optional, 0 means no age limit)
	compressOld = true
)

var (
	logWriter io.Writer
)

// Init initializes logger with lumberjack for log rotation
func Init() error {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// Configure lumberjack for log rotation
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    maxLogSize,  // megabytes
		MaxBackups: maxBackups,  // number of backups
		MaxAge:     maxAge,      // days
		Compress:   compressOld, // compress old log files
		LocalTime:  true,        // use local time for filenames
	}

	// Create multi-writer to write to both stdout and file
	logWriter = io.MultiWriter(os.Stdout, lumberjackLogger)

	// Configure standard log package
	log.SetOutput(logWriter)
	log.SetFlags(log.Ldate | log.Ltime)

	Info("Logger initialized")

	return nil
}

// Close is a no-op for lumberjack (it handles cleanup automatically)
// but kept for API compatibility
func Close() {
	Info("Logger shutting down")
}

// Helper functions for common log patterns

// getCaller returns the function name and file location of the caller
func getCaller() (string, string) {
	pc, file, line, ok := runtime.Caller(3) // 3 levels up: formatLog -> Info/Error/etc -> actual caller
	if !ok {
		return "unknown", "unknown:0"
	}
	fn := runtime.FuncForPC(pc)
	funcName := "unknown"
	if fn != nil {
		// Extract just the function name
		fullName := fn.Name()
		parts := strings.Split(fullName, ".")
		if len(parts) > 0 {
			funcName = parts[len(parts)-1]
		}
	}
	// Extract just the filename (not full path)
	fileParts := strings.Split(file, "/")
	fileName := fileParts[len(fileParts)-1]
	location := fmt.Sprintf("%s:%d", fileName, line)
	return funcName, location
}

// formatLog creates a log message with level, location, function, and message
func formatLog(level, msg string, args ...any) string {
	funcName, location := getCaller()

	// Build the log line with tabs
	logLine := fmt.Sprintf("%s\t[%s]\t[%s]\t%s", location, level, funcName, msg)

	// Add key-value pairs if provided
	if len(args) > 0 {
		for i := 0; i < len(args); i += 2 {
			if i+1 < len(args) {
				logLine += fmt.Sprintf(" %v=%v", args[i], args[i+1])
			}
		}
	}

	return logLine
}

// Info logs an informational message
func Info(msg string, args ...any) {
	log.Println(formatLog("INFO", msg, args...))
}

// Error logs an error message
func Error(msg string, args ...any) {
	log.Println(formatLog("ERROR", msg, args...))
}

// Warn logs a warning message
func Warn(msg string, args ...any) {
	log.Println(formatLog("WARN", msg, args...))
}

// Debug logs a debug message
func Debug(msg string, args ...any) {
	log.Println(formatLog("DEBUG", msg, args...))
}

// Request logs an HTTP request at the specified level
func Request(level, method, path string, status int, duration string) {
	funcName, location := getCaller()
	logLine := fmt.Sprintf("%s\t[%s]\t[%s] %s %s status=%d duration=%s",
		location, level, funcName, method, path, status, duration)
	log.Println(logLine)
}
