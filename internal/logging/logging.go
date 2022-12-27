package logging

import (
	"log"
)

type DebugLog struct {
	log        *log.Logger
	PrintDebug bool
}

var Logger *DebugLog

func (m *DebugLog) SetPrefix(prefix string) {
	m.log.SetPrefix(prefix)
}

func (m *DebugLog) Debug(args ...interface{}) {
	if m.PrintDebug {
		m.Print(args...)
	}
}

func (m *DebugLog) Debugf(format string, args ...interface{}) {
	if m.PrintDebug {
		m.Printf(format, args...)
	}
}

func (m *DebugLog) Debugln(args ...interface{}) {
	if m.PrintDebug {
		m.Println(args...)
	}
}

func (m *DebugLog) Print(args ...interface{}) {
	m.log.Print(args...)
}

func (m *DebugLog) Printf(format string, args ...interface{}) {
	m.log.Printf(format, args...)
}

func (m *DebugLog) Println(args ...interface{}) {
	m.log.Println(args...)
}

func (m *DebugLog) Fatal(args ...interface{}) {
	m.log.Fatal(args...)
}

func (m *DebugLog) Fatalf(format string, args ...interface{}) {
	m.log.Fatalf(format, args...)
}
