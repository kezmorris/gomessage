package server

import (
	"log"
)

// Logger is a logger
type Logger struct {
	logChan chan string
}

// NewLogger returns a new logger
func NewLogger() Logger {
	logChannel := make(chan string, 20)
	newLogger := Logger{
		logChannel,
	}
	return newLogger
}

// Run starts a logger goroutine
func (logger *Logger) Run() {
	go logger.logMessages()
}

// Log is where messages are given to be logged
func (logger *Logger) Log(message string) {
	logger.logChan <- message
}

func (logger *Logger) logMessages() {
	for {
		newMessage := <-logger.logChan
		log.Printf(newMessage)
	}
}
