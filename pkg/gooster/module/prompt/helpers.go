package prompt

import (
	"github.com/gdamore/tcell"
	"github.com/jumale/gooster/pkg/filesys"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

func getColorName(c tcell.Color) string {
	for name, value := range tcell.ColorNames {
		if value == c {
			return name
		}
	}
	return "black"
}

func (m *Module) check(err error, msg ...string) {
	if err == nil {
		return
	}
	if len(msg) > 0 {
		m.Log().Error(errors.WithMessage(err, msg[0]))
	} else {
		m.Log().Error(err)
	}
}

var pathRegex = regexp.MustCompile(`^(?:(?:\.{1,2}/)|(?:(?:/[^/]+)+))`)

func detectWorkDirPath(fs filesys.FileSys, command string) (path string) {
	args := strings.Split(command, " ")

	if args[0] == "cd" && len(args) >= 2 {
		path = args[1]
	} else if pathRegex.MatchString(args[0]) {
		path = args[0]
	}

	if path == "" {
		return ""
	}

	if strings.HasPrefix(path, "~") {
		ud, _ := fs.UserHomeDir()
		path = strings.Replace(path, "~", ud, 1)
	}

	if strings.HasPrefix(path, "./") {
		path = strings.Replace(path, "./", "", 1)
	}

	if !strings.HasPrefix(path, "/") {
		wd, _ := fs.Getwd()
		path = strings.Replace(wd+"/"+path, "//", "/", -1)
	}

	// return empty if it's not a directory
	info, err := fs.Stat(path)
	if err != nil || !info.IsDir() {
		return ""
	}

	return path
}

type regexList []*regexp.Regexp

func (r regexList) MatchString(s string) bool {
	for _, pattern := range r {
		if pattern.MatchString(s) {
			return true
		}
	}
	return false
}
