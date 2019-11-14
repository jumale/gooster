package testtools

import (
	"bytes"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/jumale/gooster/pkg/log"
	"github.com/pkg/errors"
)

func TestableContext() (ctx *gooster.AppContext, logs *bytes.Buffer) {
	logs = bytes.NewBuffer(nil)

	var err error
	ctx, err = gooster.NewAppContext(
		gooster.AppContextConfig{
			LogLevel:          log.Debug,
			LogFormat:         "<level>##<msg>\n",
			LogTarget:         logs,
			DelayEventManager: false,
		},
	)
	if err != nil {
		panic(errors.WithMessagef(err, "could not instantiate AppContext"))
	}

	return ctx, logs
}
