package logging

import (
	"fmt"
	"log"
	"os"
)

type DebugLog struct {
	log        log.Logger
	PrintDebug bool
}
type Level int

const (
	Debug Level = iota
	Warning
	Info
	NoLogging
)

func (l Level) String() string {
	switch l {
	case Debug:
		return "DEBUG:  "
	case Warning:
		return "WARNING:"
	case Info:
		return "INFO:   "
	default:
		return "ERROR:  "
	}
}

var Logger *DebugLog

func InitLogger() *DebugLog {
	if Logger == nil {
		Logger = &DebugLog{
			log:        *log.New(os.Stdout, "", log.Ldate|log.Ltime),
			PrintDebug: false,
		}
	}
	return Logger
}

func (m *DebugLog) SetPrefix(prefix string) {
	m.log.SetPrefix(prefix)
}

func (m *DebugLog) Println(args ...interface{}) {
	m.log.Println(args...)
}

func (m *DebugLog) Print(args ...interface{}) {
	m.log.Print(args...)
}

func (m *DebugLog) Printf(format string, args ...interface{}) {
	m.log.Printf(format, args...)
}

func (m *DebugLog) Log(l Level, args ...interface{}) {
	if l > Debug || m.PrintDebug {
		if l >= NoLogging {
			fmt.Println(args...)
		} else {
			m.Println(l, args)
		}
	}
}

func (m *DebugLog) Logf(l Level, format string, args ...interface{}) {
	if l > Debug || m.PrintDebug {
		if l >= NoLogging {
			fmt.Printf(format, args...)
		} else {
			m.Printf("%s"+format, l, args)
		}
	}
}

func (m *DebugLog) Fatal(args ...interface{}) {
	m.log.Fatal(args...)
}

func (m *DebugLog) Fatalf(format string, args ...interface{}) {
	m.log.Fatalf(format, args...)
}
