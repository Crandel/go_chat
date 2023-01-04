package logging

import (
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
	Unknown
)

func (l Level) String() string {
	switch l {
	case Debug:
		return "DEBUG"
	case Warning:
		return "WARNING"
	case Info:
		return "INFO"
	default:
		return "ERROR"
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

func (m *DebugLog) Debug(args ...interface{}) {
	if m.PrintDebug {
		m.Print(args...)
	}
}

func (m *DebugLog) Dbg(l Level, args ...interface{}) {
	m.Debug(l, ": ", args)
}

func (m *DebugLog) Debugf(format string, args ...interface{}) {
	if m.PrintDebug {
		m.Printf(format, args...)
	}
}

func (m *DebugLog) Dbgf(l Level, format string, args ...interface{}) {
	m.Debugf("%s: "+format, l, args)
}

func (m *DebugLog) Debugln(args ...interface{}) {
	if m.PrintDebug {
		m.Println(args...)
	}
}

func (m *DebugLog) Dbgln(l Level, args ...interface{}) {
	m.Debugln(l, ": ", args)
}

func (m *DebugLog) Print(args ...interface{}) {
	m.log.Print(args...)
}

func (m *DebugLog) Pr(l Level, args ...interface{}) {
	m.Print(l, ": ", args)
}

func (m *DebugLog) Printf(format string, args ...interface{}) {
	m.log.Printf(format, args...)
}

func (m *DebugLog) Prf(l Level, format string, args ...interface{}) {
	m.Printf("%s: "+format, l, args)
}

func (m *DebugLog) Println(args ...interface{}) {
	m.log.Println(args...)
}

func (m *DebugLog) Prln(l Level, args ...interface{}) {
	m.Println(l, ": ", args)
}

func (m *DebugLog) Fatal(args ...interface{}) {
	m.log.Fatal(args...)
}

func (m *DebugLog) Fatalf(format string, args ...interface{}) {
	m.log.Fatalf(format, args...)
}
