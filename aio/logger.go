package aio

import "fmt"

type Logger struct {
	IsVerbose bool
}

func (*Logger) Println(a ...interface{}) {
	fmt.Println(a...)
}

func (*Logger) Printf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

func (l *Logger) Verbose(a ...interface{}) {
	if l.IsVerbose {
		l.Println(a...)
	}
}

func (l *Logger) Verbosef(format string, a ...interface{}) {
	if l.IsVerbose {
		l.Printf(format, a...)
	}
}
