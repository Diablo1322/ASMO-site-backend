package logger

import (
	"encoding/json"
	"log"
	"time"
)

type LogLevel string

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
)

type LogEntry struct {
	Timestamp string      `json:"timestamp"`
	Level     LogLevel    `json:"level"`
	Service   string      `json:"service"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}

type Logger struct {
	service   string
	level     LogLevel
	requestID string
}

func New(service string, level LogLevel) *Logger {
	return &Logger{
		service: service,
		level:   level,
	}
}

func (l *Logger) shouldLog(level LogLevel) bool {
	levels := map[LogLevel]int{
		DEBUG: 0,
		INFO:  1,
		WARN:  2,
		ERROR: 3,
	}
	return levels[level] >= levels[l.level]
}

func (l *Logger) log(level LogLevel, message string, data interface{}) {
	if !l.shouldLog(level) {
		return
	}

	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     level,
		Service:   l.service,
		Message:   message,
		Data:      data,
		RequestID: l.requestID,
	}

	logJSON, err := json.Marshal(entry)
	if err != nil {
		log.Printf("Error marshaling log entry: %v", err)
		return
	}

	log.Println(string(logJSON))
}

func (l *Logger) Debug(message string, data interface{}) {
	l.log(DEBUG, message, data)
}

func (l *Logger) Info(message string, data interface{}) {
	l.log(INFO, message, data)
}

func (l *Logger) Warn(message string, data interface{}) {
	l.log(WARN, message, data)
}

func (l *Logger) Error(message string, data interface{}) {
	l.log(ERROR, message, data)
}

func (l *Logger) WithRequestID(requestID string) *Logger {
	return &Logger{
		service:   l.service,
		level:     l.level,
		requestID: requestID,
	}
}