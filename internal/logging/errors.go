package logging

import (
	"fmt"
	"strings"
)

type Op string

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

type CommonError struct {
	Err     error
	Op      Op
	Message string
	Level   Level
}

func (e CommonError) Error() string {
	return e.Message
}

func NewError(op Op, l Level, m string, err error) CommonError {
	er := CommonError{
		Op:      op,
		Level:   l,
		Message: m,
		Err:     err,
	}
	er.Logging()
	return er
}

func New(op Op, l Level, m string) CommonError {
	return NewError(op, l, m, nil)
}

func Tracing(e CommonError) []string {
	stack := []string{string(e.Op)}
	intError, ok := e.Err.(*CommonError)
	if !ok {
		return stack
	}
	stack = append(stack, Tracing(*intError)...)
	return stack
}
func (e *CommonError) Logging() {
	format := "%s: [%s] - %s"
	if e.Err != nil {
		format = format + fmt.Sprintf(". Error: %v", e.Err)
	}

	stack := Tracing(*e)
	finalStack := strings.Join(stack, "::")
	Logger.Debugf(format, fmt.Sprint(e.Level), finalStack, e.Message)
}
