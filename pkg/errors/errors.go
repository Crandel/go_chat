package errors

import (
	"fmt"

	"log"
)

type Op string

type Level int

const (
	Debug Level = iota
	Warning
	Info
	Unknown
)

const desiredLevel Level = Debug

func (l *Level) String() string {
	switch *l {
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

type Error struct {
	Op    Op
	Level Level
	Err   error
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: [%s] - %v", fmt.Sprint(e.Level), e.Op, e.Err.Error())
}

func New(op Op, l Level, err error) Error {
	new_err := Error{op, l, err}
	Logging(new_err, desiredLevel)
	return new_err
}

func Tracing(e *Error) []Op {
	stack := []Op{e.Op}
	intError, ok := e.Err.(*Error)
	if !ok {
		return stack
	}
	stack = append(stack, Tracing(intError)...)
	return stack
}

func Logging(e Error, desiredLevel Level) {
	const format = "%s: [%s] - %v"
	if e.Level == Unknown {
		log.Fatal(format, "ERROR", e.Op, e.Err)
	} else if e.Level > desiredLevel {
		log.Printf(format, fmt.Sprint(e.Level), e.Op, e.Err)
	}
}
