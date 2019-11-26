package testtools

import (
	"bytes"
	"fmt"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/filesys/fstub"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/jumale/gooster/pkg/log"
	"github.com/stretchr/testify/require"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

type ExtensionTester struct {
	*gooster.AppContext
	Extension    gooster.Extension
	Target       gooster.Module
	ConfigReader *ConfigReader
	Fs           *fstub.Stub
	logs         *bytes.Buffer
	assert       *require.Assertions
	events       []events.IEvent
}

func NewExtensionTester(t *testing.T, ext gooster.Extension, target gooster.Module, cfg interface{}) *ExtensionTester {
	ctx, fs, cfgReader, logs := TestableContext()
	ctx.SetCfgPath(ext.Name())
	cfgReader.ShouldReturn(ext.Name(), cfg)

	tester := &ExtensionTester{
		Extension:    ext,
		Target:       target,
		AppContext:   ctx,
		ConfigReader: cfgReader,
		Fs:           fs,
		logs:         logs,
		assert:       require.New(t),
	}

	ctx.Events().Subscribe(events.HandleWithPrio(events.AfterAllOtherChanges, func(e events.IEvent) events.IEvent {
		tester.events = append(tester.events, e)
		return e
	}))

	return tester
}

func (t *ExtensionTester) Init() error {
	return t.Extension.Init(t.Target, t.AppContext)
}

func (t *ExtensionTester) AssertInited() {
	t.assert.NoError(t.Extension.Init(t.Target, t.AppContext))
}

//func (t *ModuleTester) PressKey(key tcell.Key, r ...rune) *ModuleTester {
//	if len(r) == 0 {
//		r = append(r, 0)
//	}
//	t.module.GetInputCapture()(tcell.NewEventKey(key, r[0], tcell.ModNone))
//	return t
//}
//
func (t *ExtensionTester) SendEvent(event events.IEvent) *ExtensionTester {
	t.Events().Dispatch(event)
	return t
}

func (t *ExtensionTester) AssertFinalEvent(event events.IEvent) {
	exists := false
	for _, e := range t.events {
		if reflect.DeepEqual(e, event) {
			exists = true
			break
		}
	}

	t.assert.True(exists, "Expected event %+v has not been dispatched\nActual: %+v", event, t.events)
}

func (t *ExtensionTester) AssertHasLog(msg string, level ...log.Level) {
	if len(level) > 0 {
		msg = level[0].String() + "##" + msg
	}
	t.assertContent("App logs do not match expected lines", false, strings.Trim(t.logs.String(), "\n"), msg)
}

//
func (t *ExtensionTester) assertContent(msg string, exact bool, actual string, expectedLines ...interface{}) {
	var matchers []string
	for _, line := range expectedLines {
		switch l := line.(type) {
		case string:
			matchers = append(matchers, regexp.QuoteMeta(l))
		case *regexp.Regexp:
			matchers = append(matchers, l.String())
		default:
			t.assert.Fail("Non supported output line matcher.", "Expected string or regexp, got %T", l)
		}
	}
	expected := strings.Join(matchers, "\n")

	pattern := regexp.MustCompile(expected)
	if exact {
		pattern = regexp.MustCompile("^" + expected + "$")
	}

	if !pattern.MatchString(actual) {
		t.assert.Fail(fmt.Sprintf(
			"Not equal: \n"+
				msg+": \n"+
				"expected: \"%s\"\n"+
				"actual  : \"%s\"%s",
			strings.ReplaceAll(expected, "\n", `\n`),
			strings.ReplaceAll(actual, "\n", `\n`),
			diff(expected, actual),
		))
	}
}
