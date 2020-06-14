package utils

import (
	"os"
	"time"
)

type Logger struct {
	logPath  string
	Messages chan string
}

func NewLogger(logPath string) *Logger {
	logger := new(Logger)
	logger.logPath = logPath
	logger.Messages = make(chan string)

	return logger
}

func StartLogger(logger *Logger, quit chan bool) {
	for {
		select {
		case <-quit:
			// StartLogger quits
			return
		case msg := <-logger.Messages:
			logger.LogMessage(msg)

			// Wait 100 milliseconds between logs
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// Needs to be a routine
func (log Logger) LogMessage(msg string) {
	f, _ := os.OpenFile(log.logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	f.Write([]byte(msg + "\n"))
}
