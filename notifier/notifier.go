package notifier

import (
	"github.com/jouir/pgterminate/base"
)

// Notifier generic interface for implementing a notifier
type Notifier interface {
	Run()
	Reload()
}

// NewNotifier looks into Config to create a Console, File or Syslog notifier and pass it
// the session channel for consuming sessions structs sent by terminator
func NewNotifier(ctx *base.Context) Notifier {
	switch ctx.Config.LogDestination {
	case "file":
		return NewFile(ctx.Config.LogFile, ctx.Sessions)
	case "syslog":
		return NewSyslog(ctx.Config.SyslogFacility, ctx.Config.SyslogIdent, ctx.Sessions)
	default: // console
		return NewConsole(ctx.Sessions)
	}
}
