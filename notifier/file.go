package notifier

import (
	"github.com/jouir/pgterminate/base"
	"github.com/jouir/pgterminate/log"
	"os"
	"sync"
	"time"
)

// File structure for file notifier
type File struct {
	handle   *os.File
	name     string
	sessions chan *base.Session
	mutex    sync.Mutex
}

// NewFile creates a file notifier
func NewFile(name string, sessions chan *base.Session) Notifier {
	return &File{
		name:     name,
		sessions: sessions,
	}
}

// Run starts the file notifier
func (f *File) Run() {
	log.Info("Starting file notifier")
	f.open()
	defer f.terminate()

	for session := range f.sessions {
		timestamp := time.Now().Format(time.RFC3339)
		_, err := f.handle.WriteString(timestamp + " " + session.String() + "\n")
		base.Panic(err)
	}
}

// open opens a log file
func (f *File) open() {
	var err error
	f.handle, err = os.OpenFile(f.name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	base.Panic(err)
}

// Reload closes and re-open the file to be compatible with logrotate
func (f *File) Reload() {
	log.Info("Reloading file notifier")
	f.mutex.Lock()
	defer f.mutex.Unlock()
	log.Debugf("Re-opening log file %s\n", f.name)
	f.handle.Close()
	f.open()
}

// terminate closes the file
func (f *File) terminate() {
	log.Debugf("Closing log file %s\n", f.name)
	f.handle.Close()
}
