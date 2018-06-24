package log

import (
	"errors"
	"log"
	"os"
)

const (
	// DebugLevel for debug messages
	DebugLevel int = iota
	// InfoLevel for info messages
	InfoLevel
	// WarnLevel for warning messages
	WarnLevel
	// ErrorLevel for error messages
	ErrorLevel
	// FatalLevel for fatal messages
	FatalLevel
)

var level int

func init() {
	level = WarnLevel
}

// SetLevel configures logging level
func SetLevel(logLevel int) error {
	if logLevel < DebugLevel || logLevel > FatalLevel {
		return errors.New("Wrong logging level")
	}
	level = logLevel
	return nil
}

// Debug prints debug messages
func Debug(m string) {
	if level <= DebugLevel {
		log.Println("[DEBUG] " + m)
	}
}

// Debugf prints debug messages with printf format
func Debugf(format string, values ...interface{}) {
	if level <= DebugLevel {
		log.Printf("[DEBUG] "+format, values...)
	}
}

// Info prints info messages
func Info(m string) {
	if level <= InfoLevel {
		log.Println("[INFO] " + m)
	}
}

// Infof prints info messages with printf format
func Infof(format string, values ...interface{}) {
	if level <= InfoLevel {
		log.Printf("[INFO] "+format, values...)
	}
}

// Warn prints warning messages
func Warn(m string) {
	if level <= WarnLevel {
		log.Println("[WARN] " + m)
	}
}

// Warnf prints warning messages with printf format
func Warnf(format string, values ...interface{}) {
	if level <= WarnLevel {
		log.Printf("[WARN] "+format, values...)
	}
}

// Error prints error messages
func Error(m string) {
	if level <= ErrorLevel {
		log.Println("[ERROR] " + m)
	}
}

// Errorf prints error messages with printf format
func Errorf(format string, values ...interface{}) {
	if level <= WarnLevel {
		log.Printf("[ERROR] "+format, values...)
	}
}

// Fatal prints fatal messages and exit
func Fatal(m string) {
	if level <= FatalLevel {
		log.Println("[FATAL] " + m)
	}
	os.Exit(1)
}

// Fatalf prints fatal messages with printf format and exit
func Fatalf(format string, values ...interface{}) {
	if level <= FatalLevel {
		log.Printf("[FATAL] "+format, values...)
	}
	os.Exit(1)
}
