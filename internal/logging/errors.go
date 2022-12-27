package logging

import "fmt"

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
	format := "%s: [%s] - %s"
	if e.Err != nil {
		format = format + "err: " + e.Err.Error()
	}
	return fmt.Sprintf(format, e.Level.String(), e.Op, e.Message)
}

func NewError(op Op, l Level, m string, err error) CommonError {
	return CommonError{
		Op:      op,
		Level:   l,
		Message: m,
		Err:     err,
	}
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
