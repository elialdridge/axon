package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	debugLogger *log.Logger
	infoLogger  *log.Logger
	errorLogger *log.Logger
	logFile     *os.File
)

// Init initializes the logger
func Init() error {
	var err error
	logFile, err = os.OpenFile("axon_debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	debugLogger = log.New(logFile, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLogger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(logFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	Info("Logger initialized")
	return nil
}

// Close closes the log file
func Close() {
	if logFile != nil {
		logFile.Close()
	}
}

// Debug logs debug messages
func Debug(format string, v ...interface{}) {
	if debugLogger != nil {
		debugLogger.Printf(format, v...)
	}
}

// Info logs info messages
func Info(format string, v ...interface{}) {
	if infoLogger != nil {
		infoLogger.Printf(format, v...)
	}
}

// Error logs error messages
func Error(format string, v ...interface{}) {
	if errorLogger != nil {
		errorLogger.Printf(format, v...)
	}
}

// LogRequest logs AI request details
func LogRequest(req interface{}) {
	Debug("AI Request: %+v", req)
}

// LogResponse logs AI response details
func LogResponse(resp interface{}) {
	Debug("AI Response: %+v", resp)
}

// LogGameState logs game state details
func LogGameState(state interface{}) {
	Debug("Game State: %+v", state)
}

// LogWorldCreation logs world creation process
func LogWorldCreation(step string, data interface{}) {
	Info("World Creation - %s: %+v", step, data)
}

