package logs

import (
	"fmt"
	stdlog "log"
	"os"
)

// LogStruct describes the shape of a new log
type LogStruct struct {
	Layer   string
	Message string
}

// NewLog creates a new LogStruct
func NewLog(layer string, message string) LogStruct {
	return LogStruct{Layer: layer, Message: message}
}

// Print determines how a log should be printed and sets the timestamp
func (log *LogStruct) Print(logLevel int) string {
	levelStr := GetLogLevelString(logLevel)
	return fmt.Sprintf("[%s] (%s) %s", levelStr, log.Layer, log.Message)
}

// Logger describes the methods available on a logger
type Logger interface {
	SetLevel(int)
	Fatal(LogStruct)
	Error(LogStruct)
	Warning(LogStruct)
	Info(LogStruct)
	Debug(LogStruct)
	Verbose(LogStruct)
}

const (
	// LogLevelFatal prints only messages which cause critical application failure
	LogLevelFatal = iota
	// LogLevelError prints only serious application errors
	LogLevelError = iota
	// LogLevelWarning prints potential issues in application
	LogLevelWarning = iota
	// LogLevelInfo prints general information of application working as expected
	LogLevelInfo = iota
	// LogLevelDebug prints detailed information intended for debugging issues
	LogLevelDebug = iota
	// LogLevelVerbose prints all and every detail
	LogLevelVerbose = iota
)

// GetLogLevelString converts iota into string representation
func GetLogLevelString(level int) string {
	switch level {
	case LogLevelFatal:
		return "FATAL"
	case LogLevelError:
		return "ERROR"
	case LogLevelWarning:
		return "WARN"
	case LogLevelInfo:
		return "INFO"
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelVerbose:
		return "VERBOSE"
	}
	return "UNKNOWN"
}

// StdLogger implements Logger using the standard library
type StdLogger struct {
	level   int
	fatal   *stdlog.Logger
	error   *stdlog.Logger
	warning *stdlog.Logger
	info    *stdlog.Logger
	debug   *stdlog.Logger
	verbose *stdlog.Logger
}

// NewStdLogger creates a new StdLogger instance
func NewStdLogger(level int) *StdLogger {
	var (
		fatal   = stdlog.New(os.Stderr, "", stdlog.Ldate|stdlog.Ltime)
		error   = stdlog.New(os.Stderr, "", stdlog.Ldate|stdlog.Ltime)
		warning = stdlog.New(os.Stdout, "", stdlog.Ldate|stdlog.Ltime)
		info    = stdlog.New(os.Stdout, "", stdlog.Ldate|stdlog.Ltime)
		debug   = stdlog.New(os.Stdout, "", stdlog.Ldate|stdlog.Ltime)
		verbose = stdlog.New(os.Stdout, "", stdlog.Ldate|stdlog.Ltime)
	)

	return &StdLogger{level, fatal, error, warning, info, debug, verbose}
}

// SetLevel updates logger level
func (logger *StdLogger) SetLevel(level int) {
	logger.level = level
}

// Fatal prints fatal logs
func (logger *StdLogger) Fatal(log LogStruct) {
	if logger.level >= LogLevelFatal {
		logger.fatal.Println(log.Print(LogLevelFatal))
	}
}

// Error prints fatal logs
func (logger *StdLogger) Error(log LogStruct) {
	if logger.level >= LogLevelError {
		logger.error.Println(log.Print(LogLevelError))
	}
}

// Warning prints fatal logs
func (logger *StdLogger) Warning(log LogStruct) {
	println("warning called msg: %s. level: %s", log.Message, GetLogLevelString(LogLevelWarning))
	if logger.level >= LogLevelWarning {
		logger.warning.Println(log.Print(LogLevelWarning))
	}
}

// Info prints fatal logs
func (logger *StdLogger) Info(log LogStruct) {
	if logger.level >= LogLevelInfo {
		logger.info.Println(log.Print(LogLevelInfo))
	}
}

// Debug prints fatal logs
func (logger *StdLogger) Debug(log LogStruct) {
	if logger.level >= LogLevelDebug {
		logger.debug.Println(log.Print(LogLevelDebug))
	}
}

// Verbose prints fatal logs
func (logger *StdLogger) Verbose(log LogStruct) {
	if logger.level >= LogLevelVerbose {
		logger.verbose.Println(log.Print(LogLevelVerbose))
	}
}
