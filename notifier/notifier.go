package notifier

import (
	"github.com/jouir/pgterminate/base"
)

// Notifier generic interface for implementing a notifier
type Notifier interface {
	Run()
	Reload()
}

// NewNotifier looks into Config to create a File or Console notifier and pass it
// the session channel for consuming sessions structs sent by terminator
func NewNotifier(ctx *base.Context) Notifier {
	if ctx.Config.LogFile != "" {
		return NewFile(ctx.Config.LogFile, ctx.Sessions)
	}
	if ctx.Config.SyslogFacility != "" {
		return NewSyslog(ctx.Config.SyslogFacility, ctx.Config.SyslogIdent, ctx.Sessions)
	}
	return NewConsole(ctx.Sessions)
}
