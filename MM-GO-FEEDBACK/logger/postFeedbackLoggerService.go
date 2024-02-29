package MMLogger

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func (l *Logger) log(level LogLevel, logMessage interface{}, args ...interface{}) {
	if level > l.level {
		return
	}
	var message string
	switch m := logMessage.(type) {
	case string:
		message = fmt.Sprintf(m, args...)
	case fmt.Stringer:
		message = m.String()
	default:
		v := reflect.ValueOf(m)
		if v.Kind() == reflect.Struct {
			var b strings.Builder
			b.WriteString("{")
			for i := 0; i < v.NumField(); i++ {
				if i > 0 {
					b.WriteString(", ")
				}
				fmt.Fprintf(&b, "%s: %v", v.Type().Field(i).Name, v.Field(i))
			}
			b.WriteString("}")
			message = b.String()
		} else {
			message = fmt.Sprintf("%+v", m)
		}
	}
	message = l.maskSensitiveWords(message)
	if len(message) > maxMessageLength {
		message = message[:maxMessageLength] + " [log trimmed]"
	}
	timestamp := time.Now().Format("2024-02-28 15:04:05")
	logLevel := [...]string{"INFO", "DEBUG", "ERROR"}[level]
	logMessageFmt := fmt.Sprintf("TimeStamp: [%s] LogLevel: [%s] LogMessage: %s\n", timestamp, logLevel, message)
	fmt.Print(logMessageFmt)
}

func (l *Logger) Info(format interface{}, args ...interface{}) {
	l.log(INFO, format, args...)
}
func (l *Logger) Debug(format interface{}, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

func (l *Logger) Error(format interface{}, args ...interface{}) {
	l.log(ERROR, format, args...)
}

///

const maxMessageLength = 1000

type LogLevel int64

const (
	INFO LogLevel = iota
	DEBUG
	ERROR
)

type Logger struct {
	level       LogLevel
	wordsToMask []string
}

func NewLogger() *Logger {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	var logStr = os.Getenv("Log_Level")
	logInt, _ := strconv.ParseInt(logStr, 10, 64)
	var loglevel = LogLevel(logInt)
	// level:= loglevel
	switch loglevel {
	case INFO:
		loglevel = INFO

	case DEBUG:
		loglevel = DEBUG

	case ERROR:
		loglevel = ERROR

	default:
		loglevel = INFO
	}
	wordsToMask := []string{"Hello", "World", "PlanCost"}
	return &Logger{
		level:       loglevel,
		wordsToMask: wordsToMask,
	}

}

func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

//

func (l *Logger) maskSensitiveWords(message string) string {
	maskedMessage := message
	for _, word := range l.wordsToMask {
		maskedMessage = strings.Replace(maskedMessage, word, strings.Repeat("*", len(word)), -1)
	}
	return maskedMessage
}
