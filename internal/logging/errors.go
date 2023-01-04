package logging

import (
	"fmt"
	"strings"
)

type Stk string

type CommonError struct {
	Err     error
	Stk     Stk
	Message string
}

func (e CommonError) Error() string {
	return e.Message
}

func NewError(op Stk, m string, err error) CommonError {
	er := CommonError{
		Stk:     op,
		Message: m,
		Err:     err,
	}
	er.Logging()
	return er
}

func New(op Stk, m string) CommonError {
	return NewError(op, m, nil)
}

func Tracing(e CommonError) []string {
	stack := []string{string(e.Stk)}
	intError, ok := e.Err.(*CommonError)
	if !ok {
		return stack
	}
	stack = append(stack, Tracing(*intError)...)
	return stack
}
func (e *CommonError) Logging() {
	format := "[%s] - %s"
	if e.Err != nil {
		format = format + fmt.Sprintf(". Error: %v", e.Err)
	}

	stack := Tracing(*e)
	finalStack := strings.Join(stack, "::")
	Logger.Debugf(format, finalStack, e.Message)
}
