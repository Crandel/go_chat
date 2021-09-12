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

type CommonError struct {
	Op      Op
	Level   Level
	Message string
	Err     error
}

func (e CommonError) Error() string {
	return e.Message
}

func NewError(op Op, l Level, m string, err error) CommonError {
	new_err := CommonError{op, l, m, err}
	Logging(new_err, desiredLevel)
	return new_err

}

func New(op Op, l Level, m string) CommonError {
	return NewError(op, l, m, nil)
}

func Tracing(e *CommonError) []Op {
	stack := []Op{e.Op}
	intError, ok := e.Err.(*CommonError)
	if !ok {
		return stack
	}
	stack = append(stack, Tracing(intError)...)
	return stack
}

func Logging(e CommonError, desiredLevel Level) {
	format := "%s: [%s] - %s"
	if e.Err != nil {
		format = format + fmt.Sprintf(". Error: %v", e.Err)
	}
	if e.Level == Unknown {
		log.Fatal(format, "ERROR", e.Op, e.Message)
	} else if e.Level > desiredLevel {
		log.Printf(format, fmt.Sprint(e.Level), e.Op, e.Message)
	}
}
