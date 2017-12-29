package log

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/dhickie/openhub/log/models"
)

// Logger is the logger object which stores logs before being persisted to disk
type Logger struct {
	logs    []models.Log
	logLock *sync.Mutex
	logFile *os.File
}

// NewLogger creates a new logger pointing to the specified file
func NewLogger(fileName string) (Logger, error) {
	logs := make([]models.Log, 0)
	logFile, err := openFile(fileName)
	if err != nil {
		return Logger{}, err
	}

	return Logger{
		logs:    logs,
		logLock: &sync.Mutex{},
		logFile: logFile,
	}, nil
}

// Log adds the specified log object to the logger
func (l *Logger) Log(log models.Log) {
	l.logLock.Lock()
	l.logs = append(l.logs, log)
	l.logLock.Unlock()
}

// Flush flushes all the logs out to the log file
func (l *Logger) Flush() {
	// Lock the logs, then write them out line by line
	l.logLock.Lock()
	for _, v := range l.logs {
		logType := strings.ToUpper(v.Type)
		l.logFile.WriteString(fmt.Sprintf("%v - %v - %v\r\n", logType, v.Timestamp, v.Message))
	}
	logCount := len(l.logs)
	l.logs = l.logs[logCount:]
	l.logLock.Unlock()
}

// NewFile moves the logger on to a new log file
func (l *Logger) NewFile(fileName string) error {
	// Log the log lock so that nobody can flush while we're switching files
	l.logLock.Lock()
	defer l.logLock.Unlock()

	// Close the existing file
	err := l.logFile.Close()
	if err != nil {
		return err
	}

	// Open up the new one
	file, err := openFile(fileName)
	if err != nil {
		return err
	}

	l.logFile = file
	return nil
}

func openFile(fileName string) (*os.File, error) {
	return os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}
