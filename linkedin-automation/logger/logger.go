package logger

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// Logger wraps logrus logger with custom methods
type Logger struct {
	*logrus.Logger
}

// InitLogger initializes and returns a new logger instance
func InitLogger() *Logger {
	log := logrus.New()

	// Set log level from environment
	level := os.Getenv("LOG_LEVEL")
	switch level {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}

	// Set formatter
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	})

	// Set output to file if specified
	logFile := os.Getenv("LOG_FILE")
	if logFile != "" {
		// Create log directory if it doesn't exist
		logDir := filepath.Dir(logFile)
		if err := os.MkdirAll(logDir, 0755); err == nil {
			file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err == nil {
				log.SetOutput(file)
			}
		}
	}

	return &Logger{log}
}

// WithField adds a single field to the logger
func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.Logger.WithField(key, value)
}

// WithFields adds multiple fields to the logger
func (l *Logger) WithFields(fields map[string]interface{}) *logrus.Entry {
	return l.Logger.WithFields(logrus.Fields(fields))
}

// Action logs an action being performed
func (l *Logger) Action(action string, details map[string]interface{}) {
	l.WithFields(details).Infof("🎯 Action: %s", action)
}

// Success logs a successful operation
func (l *Logger) Success(message string, details map[string]interface{}) {
	l.WithFields(details).Infof("✅ %s", message)
}

// Failure logs a failed operation
func (l *Logger) Failure(message string, err error, details map[string]interface{}) {
	fields := make(map[string]interface{})
	for k, v := range details {
		fields[k] = v
	}
	if err != nil {
		fields["error"] = err.Error()
	}
	l.WithFields(fields).Errorf("❌ %s", message)
}

// Stealth logs stealth-related actions
func (l *Logger) Stealth(action string, details map[string]interface{}) {
	l.WithFields(details).Debugf("🥷 Stealth: %s", action)
}
