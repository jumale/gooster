package gooster

import (
	"fmt"
	"github.com/jumale/gooster/pkg/events"
)

type output struct {
	em events.Manager
}

func (o *output) Write(p []byte) (n int, err error) {
	o.em.Dispatch(EventOutput{Data: p})
	return len(p), nil
}

func (o *output) WriteString(s string) {
	o.em.Dispatch(EventOutput{Data: []byte(s)})
}

func (o *output) WriteBytes(b []byte) {
	o.em.Dispatch(EventOutput{Data: b})
}

func (o *output) WriteLine(s string) {
	o.em.Dispatch(EventOutput{Data: []byte(s + "\n")})
}

func (o *output) WriteBytesLine(b []byte) {
	o.em.Dispatch(EventOutput{Data: append(b, 10)})
}

func (o *output) WriteF(format string, a ...interface{}) {
	o.em.Dispatch(EventOutput{
		Data: []byte(fmt.Sprintf(format, a...)),
	})
}

func (o *output) WriteBytesF(format []byte, a ...interface{}) {
	o.WriteF(string(format), a...)
}
