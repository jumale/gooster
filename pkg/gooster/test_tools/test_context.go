package testtools

import (
	"bytes"
	"github.com/jumale/gooster/pkg/filesys/fstub"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/jumale/gooster/pkg/log"
	"github.com/pkg/errors"
)

func TestableContext() (ctx *gooster.AppContext, fs *fstub.Stub, cfg *ConfigReader, logs *bytes.Buffer) {
	logs = bytes.NewBuffer(nil)
	cfg = &ConfigReader{stubs: make(map[string]interface{})}
	fs = fstub.New(fstub.Config{
		WorkDir: "/current",
		HomeDir: "/home",
	})

	var err error
	ctx, err = gooster.NewAppContext(
		gooster.AppContextConfig{
			LogLevel:          log.Debug,
			LogFormat:         "<level>##<msg>\n",
			LogTarget:         logs,
			DelayEventManager: false,
			ConfigReader:      cfg,
			FileSys:           fs,
		},
	)
	if err != nil {
		panic(errors.WithMessagef(err, "could not instantiate AppContext"))
	}

	return ctx, fs, cfg, logs
}
