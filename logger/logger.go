package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

// Log levels
const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	FATAL
)

type LOGLEVEL int

// global logger
var (
	logger   *log.Logger
	logLevel LOGLEVEL
)

// getLogger initializes the logger if it is not already initialized and returns it.
func getLogger() *log.Logger {
	if logger == nil {
		filename := "/tmp/kamaizen.log"
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		logger = log.New(file, "[KamaiZen] ", log.Ldate|log.Ltime)
	}
	return logger
}

// SetLogLevel sets the current log level.
func SetLogLevel(level LOGLEVEL) {
	logLevel = level
}

func writeLog(level string, v ...interface{}) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = -1
	}
	getLogger().Println(append([]interface{}{fmt.Sprintf("%s:%d [%s]", file, line, level)}, v...)...)
}

// Info logs an info message if the current log level is INFO or lower.
func Info(v ...interface{}) {
	if logLevel <= INFO {
		writeLog("INFO", v...)
	}
}

// Infof logs a formatted info message if the current log level is INFO or lower.
func Infof(format string, v ...interface{}) {
	if logLevel <= INFO {
		writeLog("INFO", fmt.Sprintf(format, v...))
	}
}

// Debug logs a debug message if the current log level is DEBUG.
func Debug(v ...interface{}) {
	if logLevel <= DEBUG {
		writeLog("DEBUG", v...)
	}
}

// Debugf logs a formatted debug message if the current log level is DEBUG.
func Debugf(format string, v ...interface{}) {
	if logLevel <= DEBUG {
		writeLog("DEBUG", fmt.Sprintf(format, v...))
	}
}

func Warn(v ...interface{}) {
	if logLevel <= WARN {
		writeLog("WARN", v...)
	}
}

// Warnf logs a formatted warning message if the current log level is WARN or lower.
func Warnf(format string, v ...interface{}) {
	if logLevel <= WARN {
		writeLog("WARN", fmt.Sprintf(format, v...))
	}
}

// Error logs an error message if the current log level is ERROR or lower.
func Error(v ...interface{}) {
	if logLevel <= ERROR {
		writeLog("ERROR", v...)
	}
}

// Errorf logs a formatted error message if the current log level is ERROR or lower.
func Errorf(format string, v ...interface{}) {
	if logLevel <= ERROR {
		writeLog("ERROR", fmt.Sprintf(format, v...))
	}
}

// Fatal logs a fatal message and exits the program.
func Fatal(v ...interface{}) {
	writeLog("FATAL", v...)
}

// Fatalf logs a formatted fatal message and exits the program.
func Fatalf(format string, v ...interface{}) {
	writeLog("FATAL", fmt.Sprintf(format, v...))
}
