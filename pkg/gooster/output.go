package gooster

import (
	"fmt"
	"github.com/jumale/gooster/pkg/events"
)

type output struct {
	em events.Manager
}

func (o *output) WriteBytes(b []byte) {
	o.em.Dispatch(EventOutput{Data: b})
}

func (o *output) Write(p []byte) (n int, err error) {
	o.WriteBytes(p)
	return len(p), nil
}

func (o *output) WriteString(s string) {
	o.WriteBytes([]byte(s))
}

func (o *output) WriteLine(s string) {
	o.WriteBytesLine([]byte(s))
}

func (o *output) WriteBytesLine(b []byte) {
	o.WriteBytes(append(b, 10))
}

func (o *output) WriteF(format string, a ...interface{}) {
	o.WriteString(fmt.Sprintf(format, a...))
}

func (o *output) WriteBytesF(format []byte, a ...interface{}) {
	o.WriteF(string(format), a...)
}
