package history

import (
	"bufio"
	"github.com/jumale/gooster/pkg/filesys"
	"github.com/jumale/gooster/pkg/log"
	"github.com/pkg/errors"
	"os"
	"strings"
)

type Config struct {
	HistoryFile string
	Log         log.Logger
	FileSys     filesys.FileSys
}

type Manager struct {
	set      map[string]struct{}
	stack    []string
	index    int
	filePath string
	log      log.Logger
	fs       filesys.FileSys
}

func NewManager(cfg Config) *Manager {
	if cfg.FileSys == nil {
		cfg.FileSys = filesys.Default{}
	}

	mng := &Manager{
		set: make(map[string]struct{}),
		log: log.EmptyLogger{},
		fs:  cfg.FileSys,
	}
	if cfg.Log != nil {
		mng.log = cfg.Log
	}

	if strings.HasPrefix(cfg.HistoryFile, "~") {
		if dir, err := mng.fs.UserHomeDir(); err == nil {
			cfg.HistoryFile = strings.Replace(cfg.HistoryFile, "~", dir, 1)
		}
	}
	mng.filePath = cfg.HistoryFile
	if mng.filePath != "" {
		mng.loadHistoryLines(mng.filePath)
	}

	mng.Reset()

	return mng
}

func (h *Manager) Add(cmd string) {
	h.log.DebugF("History add: `%s`", cmd)
	h.Reset()
	h.add(cmd)
	h.write(cmd)
}

func (h *Manager) add(cmd string) {
	if _, exists := h.set[cmd]; !exists {
		h.set[cmd] = struct{}{}
		h.stack = append(h.stack, cmd)
		return
	}

	for idx, item := range h.stack {
		if item == cmd {
			h.stack = append(h.stack[:idx], h.stack[idx+1:]...)
			h.stack = append(h.stack, cmd)
			return
		}
	}
}

func (h *Manager) Reset() {
	h.index = -1
}

func (h *Manager) Prev() string {
	ln := len(h.stack)
	if ln == 0 {
		h.log.Debug("history: there is no prev")
		return ""
	}

	if !h.IsActive() {
		h.index = ln
	}
	h.index--
	// loop back if reached end of list
	if h.index < 0 {
		h.index = ln - 1
	}

	return h.Current()
}

func (h *Manager) Next() string {
	ln := len(h.stack)
	if !h.IsActive() {
		h.log.Debug("history: list is not active")
		return ""
	}

	h.index++
	if h.index >= ln {
		h.Reset()
		h.log.Debug("history: there is no next")
		return ""
	}

	return h.Current()
}

func (h *Manager) Current() string {
	return h.stack[h.index]
}

func (h *Manager) Index() int {
	return h.index
}

func (h *Manager) IsFirst() bool {
	return h.index == 0
}

func (h *Manager) IsLast() bool {
	return h.index == len(h.stack)-1
}

func (h *Manager) IsActive() bool {
	return h.index != -1
}

func (h *Manager) Find(search string) string {
	for _, item := range h.stack {
		if strings.Contains(item, search) {
			return item
		}
	}
	return ""
}

func (h *Manager) Filter(search string) []string {
	var found []string
	for _, item := range h.stack {
		if strings.Contains(item, search) {
			found = append(found, item)
		}
	}
	return found
}

func (h *Manager) loadHistoryLines(filePath string) *Manager {
	f, err := h.fs.Open(filePath)
	if err != nil {
		h.log.Error(errors.WithMessage(err, "loading bash history file"))
	}
	defer func() {
		if err := f.Close(); err != nil {
			h.log.Error(errors.WithMessage(err, "closing bash history file"))
		}
	}()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		h.add(sc.Text())
	}

	return h
}

func (h *Manager) write(cmd string) {
	f, err := h.fs.OpenFile(h.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		h.log.Error(errors.WithMessage(err, "opening hash history file"))
	}
	defer func() {
		if err := f.Close(); err != nil {
			h.log.Error(errors.WithMessage(err, "closing bash history file"))
		}
	}()

	if _, err := f.Write([]byte("\n" + cmd)); err != nil {
		h.log.Error(errors.WithMessage(err, "writing to bash history file"))
	}
}
