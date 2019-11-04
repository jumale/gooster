package testtools

import (
	"bytes"
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/events"
	"github.com/jumale/gooster/pkg/gooster"
	"github.com/jumale/gooster/pkg/log"
	"github.com/pkg/errors"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/stretchr/testify/assert"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestableModule(t *testing.T, m gooster.Module) *ModuleTester {
	logs := bytes.NewBuffer(nil)
	ctx, err := gooster.NewAppContext(
		gooster.AppContextConfig{
			LogLevel:          log.Debug,
			LogFormat:         "<level>##<msg>\n",
			LogTarget:         logs,
			DelayEventManager: false,
		},
	)
	if err != nil {
		panic(errors.WithMessagef(err, "could not create a module tester for %T", m))
	}

	err = m.Init(ctx)
	if err != nil {
		panic(errors.WithMessagef(err, "could not init module %T", m))
	}

	tester := &ModuleTester{
		AppContext: ctx,
		module:     m,
		screen:     NewScreenStub(10, 10),
		output:     bytes.NewBuffer(nil),
		logs:       logs,
		assert:     assert.New(t),
	}

	ctx.Events().Subscribe(events.HandleWithPrio(-1000000, func(e events.IEvent) events.IEvent {
		switch event := e.(type) {
		case gooster.EventOutput:
			tester.output.Write(event.Data)
		case gooster.EventDraw:
			tester.draw()
		}
		return e
	}))

	return tester
}

type ModuleTester struct {
	*gooster.AppContext
	module gooster.Module
	screen *screenStub
	output *bytes.Buffer
	logs   *bytes.Buffer
	assert *assert.Assertions
}

func (t *ModuleTester) SetSize(width, height int) *ModuleTester {
	t.screen = NewScreenStub(width, height)
	return t
}

func (t *ModuleTester) Draw() *ModuleTester {
	t.draw()
	return t
}

func (t *ModuleTester) PressKey(key tcell.Key, r ...rune) *ModuleTester {
	if len(r) == 0 {
		r = append(r, 0)
	}
	t.module.GetInputCapture()(tcell.NewEventKey(key, r[0], tcell.ModNone))
	return t
}

func (t *ModuleTester) SendEvent(event events.IEvent) *ModuleTester {
	t.Events().Dispatch(event)
	return t
}

func (t *ModuleTester) AssertView(expectedLines ...interface{}) {
	lineWidth := strconv.Itoa(t.screen.width)
	for i, line := range expectedLines {
		switch l := line.(type) {
		case string:
			if len(l) < t.screen.width {
				expectedLines[i] = fmt.Sprintf("%-"+lineWidth+"s", l)
			}
		case *regexp.Regexp:
			expectedLines[i] = regexp.MustCompile(l.String() + `\s*`)
		}
	}

	t.assertContent("Module view does not match expected lines", true, t.View(), expectedLines...)
}

func (t *ModuleTester) AssertViewHasLines(expectedLines ...interface{}) {
	t.assertContent("Module view does not match expected lines", false, t.View(), expectedLines...)
}

func (t *ModuleTester) AssertOutputHasLines(expectedLines ...interface{}) {
	t.assertContent("App output does not match expected lines", false, strings.Trim(t.output.String(), "\n"), expectedLines...)
}

func (t *ModuleTester) AssertHasLog(msg string, level ...log.Level) {
	if len(level) > 0 {
		msg = level[0].String() + "##" + msg
	}
	t.assertContent("App logs do not match expected lines", false, strings.Trim(t.logs.String(), "\n"), msg)
}

func (t *ModuleTester) assertContent(msg string, exact bool, actual string, expectedLines ...interface{}) {
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

func diff(expected, actual string) string {
	diff, _ := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:        difflib.SplitLines(expected),
		B:        difflib.SplitLines(actual),
		FromFile: "Expected",
		FromDate: "",
		ToFile:   "Actual",
		ToDate:   "",
		Context:  1,
	})
	return "\n\nDiff:\n" + diff
}

func (t *ModuleTester) View() string {
	return t.screen.GetView()
}

func (t *ModuleTester) draw() {
	t.module.Draw(t.screen)
}
