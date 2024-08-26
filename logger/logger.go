package logger

import (
	"log"
	"os"
)

// Log levels
const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	FATAL
)

// global logger
var (
	logger   *log.Logger
	logLevel int
)

// GetLogger initializes the logger if it is not already initialized and returns it.
func GetLogger() *log.Logger {
	if logger == nil {
		filename := "/tmp/kamaizen.log"
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		logger = log.New(file, "[KamaiZen] ", log.Ldate|log.Ltime|log.Lshortfile)
	}
	return logger
}

// SetLogLevel sets the current log level.
func SetLogLevel(level int) {
	logLevel = level
}

// Info logs an info message if the current log level is INFO or lower.
func Info(v ...interface{}) {
	if logLevel <= INFO {
		GetLogger().Println(append([]interface{}{"[INFO]"}, v...)...)
	}
}

// Infof logs a formatted info message if the current log level is INFO or lower.
func Infof(format string, v ...interface{}) {
	if logLevel <= INFO {
		GetLogger().Printf("[INFO] "+format, v...)
	}
}

// Debug logs a debug message if the current log level is DEBUG.
func Debug(v ...interface{}) {
	if logLevel <= DEBUG {
		GetLogger().Println(append([]interface{}{"[DEBUG]"}, v...)...)
	}
}

// Debugf logs a formatted debug message if the current log level is DEBUG.
func Debugf(format string, v ...interface{}) {
	if logLevel <= DEBUG {
		GetLogger().Printf("[DEBUG] "+format, v...)
	}
}

func Warn(v ...interface{}) {
	if logLevel <= WARN {
		GetLogger().Println(append([]interface{}{"[WARN]"}, v...)...)
	}
}

// Warnf logs a formatted warning message if the current log level is WARN or lower.
func Warnf(format string, v ...interface{}) {
	if logLevel <= WARN {
		GetLogger().Printf("[WARN] "+format, v...)
	}
}

// Error logs an error message if the current log level is ERROR or lower.
func Error(v ...interface{}) {
	if logLevel <= ERROR {
		GetLogger().Println(append([]interface{}{"[ERROR]"}, v...)...)
	}
}

// Errorf logs a formatted error message if the current log level is ERROR or lower.
func Errorf(format string, v ...interface{}) {
	if logLevel <= ERROR {
		GetLogger().Printf("[ERROR] "+format, v...)
	}
}

// Fatal logs a fatal message and exits the program.
func Fatal(v ...interface{}) {
	GetLogger().Println(append([]interface{}{"[FATAL]"}, v...)...)
}

// Fatalf logs a formatted fatal message and exits the program.
func Fatalf(format string, v ...interface{}) {
	GetLogger().Fatalf("[FATAL] "+format, v...)
}
