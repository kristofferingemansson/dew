package dew

import (
	"fmt"
	"runtime"
)

type domainErrorWrapper struct {
	domain    error
	err       error
	msg       *string
	callStack callStack
}

type callStack []uintptr

// StackFrame represents a single line in the stack trace
type StackFrame string

// Trace represents a complete stack trace
type Trace []StackFrame

// New creates new Domain error wrapper error
func New(domain error, v ...interface{}) error {
	err := &domainErrorWrapper{
		domain:    domain,
		callStack: getCallStack(),
	}

	for i, w := range v {
		if e, ok := w.(error); ok && i == 0 {
			err.err = e
			continue
		}
		if m, ok := w.(string); ok {
			ms := fmt.Sprintf(m, v[i+1:]...)
			err.msg = &ms
		}
		break
	}

	return err
}

// Error formats error as string
func (e *domainErrorWrapper) Error() string {
	var ret string
	if e.domain != nil {
		ret = e.domain.Error()
	}
	if e.msg != nil {
		if ret != "" {
			ret += ": "
		}
		ret += *e.msg
	}
	if e.err != nil {
		if ret != "" {
			ret += ": "
		}
		ret += e.err.Error()
	}
	return ret
}

// Domain returns last domain error
func Domain(err error) error {
	if e, ok := err.(*domainErrorWrapper); ok {
		return e.domain
	}
	return err
}

// Cause returnes first original error
func Cause(err error) error {
	if e, ok := err.(*domainErrorWrapper); ok {
		if e.err != nil {
			return Cause(e.err)
		}
		return e.domain
	}
	return err
}

// StackTrace returns stacktrace of error as string array
func StackTrace(err error) Trace {
	ret := make(Trace, 0)
	if e, ok := err.(*domainErrorWrapper); ok {
		frames := runtime.CallersFrames(e.callStack)
		for {
			frame, more := frames.Next()
			stackFrame := fmt.Sprintf("%v() @ %v:%v", frame.Function, frame.File, frame.Line)
			ret = append(ret, StackFrame(stackFrame))
			if !more {
				break
			}
		}
	}
	return ret
}

// getCallStack retrieves current callstack
func getCallStack() callStack {
	buf := make(callStack, 32)
	n := runtime.Callers(3, buf)
	return buf[0:n]
}
