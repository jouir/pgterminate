package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"sync"
	"syscall"

	"github.com/jouir/pgterminate/base"
	"github.com/jouir/pgterminate/log"
	"github.com/jouir/pgterminate/notifier"
	"github.com/jouir/pgterminate/terminator"
	"golang.org/x/crypto/ssh/terminal"
)

// AppVersion stores application version at compilation time
var AppVersion string

// AppName to store application name
var AppName string = "pgterminate"

// GitCommit to set git commit at compilation time (can be empty)
var GitCommit string

// GoVersion to set Go version at compilation time
var GoVersion string

func main() {
	var err error
	config := base.NewConfig()

	quiet := flag.Bool("quiet", false, "Quiet mode")
	verbose := flag.Bool("verbose", false, "Verbose mode")
	debug := flag.Bool("debug", false, "Debug mode")
	version := flag.Bool("version", false, "Print version")
	flag.StringVar(&config.File, "config", "", "Configuration file")
	flag.StringVar(&config.Host, "host", "", "Instance host address")
	flag.IntVar(&config.Port, "port", 0, "Instance port")
	flag.StringVar(&config.User, "user", "", "Instance username")
	flag.StringVar(&config.Password, "password", "", "Instance password")
	flag.StringVar(&config.Database, "database", "", "Instance database")
	prompt := flag.Bool("prompt-password", false, "Prompt for password")
	flag.Float64Var(&config.Interval, "interval", 1, "Time to sleep between iterations in seconds")
	flag.IntVar(&config.ConnectTimeout, "connect-timeout", 3, "Connection timeout in seconds")
	flag.Float64Var(&config.IdleTimeout, "idle-timeout", 0, "Time for idle connections to be terminated in seconds")
	flag.Float64Var(&config.ActiveTimeout, "active-timeout", 0, "Time for active connections to be terminated in seconds")
	flag.StringVar(&config.LogDestination, "log-destination", "console", "Log destination between 'console', 'syslog' or 'file'")
	flag.StringVar(&config.LogFile, "log-file", "", "Write logs to a file")
	flag.StringVar(&config.LogFormat, "log-format", "pid=%p user=%u db=%d client=%r state=%s state_duration=%m query=%q", "Represent messages using this format")
	flag.StringVar(&config.PidFile, "pid-file", "", "Write process id into a file")
	flag.StringVar(&config.SyslogIdent, "syslog-ident", "pgterminate", "Define syslog tag")
	flag.StringVar(&config.SyslogFacility, "syslog-facility", "", "Define syslog facility from LOCAL0 to LOCAL7")
	flag.Var(&config.IncludeUsers, "include-user", "Terminate only this user (can be called multiple times)")
	flag.StringVar(&config.IncludeUsersRegex, "include-users-regex", "", "Terminate users matching this regexp")
	flag.Var(&config.ExcludeUsers, "exclude-user", "Ignore this user (can be called multiple times)")
	flag.StringVar(&config.ExcludeUsersRegex, "exclude-users-regex", "", "Ignore users matching this regexp")
	flag.Var(&config.IncludeDatabases, "include-database", "Terminate only this database (can be called multiple times)")
	flag.StringVar(&config.IncludeDatabasesRegex, "include-databases-regex", "", "Terminate databases matching this regexp")
	flag.Var(&config.ExcludeDatabases, "exclude-database", "Ignore this database (can be called multiple times)")
	flag.StringVar(&config.ExcludeDatabasesRegex, "exclude-databases-regex", "", "Ignore databases matching this regexp")
	flag.BoolVar(&config.ExcludeListeners, "exclude-listeners", false, "Ignore sessions listening for events")
	flag.BoolVar(&config.Cancel, "cancel", false, "Cancel sessions instead of terminate")
	flag.Parse()

	log.SetLevel(log.WarnLevel)
	if *debug {
		log.SetLevel(log.DebugLevel)
	}
	if *verbose {
		log.SetLevel(log.InfoLevel)
	}
	if *quiet {
		log.SetLevel(log.ErrorLevel)
	}

	if *version {
		if AppVersion == "" {
			AppVersion = "unknown"
		}
		showVersion()
		return
	}

	if *prompt {
		fmt.Print("Password:")
		bytes, err := terminal.ReadPassword(syscall.Stdin)
		base.Panic(err)
		config.Password = string(bytes)
		fmt.Print("\n")
	}

	if config.File != "" {
		err = config.Read(config.File)
		base.Panic(err)
	}

	if config.ActiveTimeout == 0 && config.IdleTimeout == 0 {
		log.Fatal("Parameter -active-timeout or -idle-timeout required")
	}

	if config.LogDestination != "console" && config.LogDestination != "file" && config.LogDestination != "syslog" {
		log.Fatal("Log destination must be 'console', 'file' or 'syslog'")
	}

	if config.LogDestination == "syslog" && config.SyslogFacility != "" {
		matched, err := regexp.MatchString("^LOCAL[0-7]$", config.SyslogFacility)
		base.Panic(err)
		if !matched {
			log.Fatal("Syslog facility must range from LOCAL0 to LOCAL7")
		}
	}

	err = config.CompileRegexes()
	base.Panic(err)
	config.CompileFilters()

	if config.PidFile != "" {
		writePid(config.PidFile)
		defer removePid(config.PidFile)
	}

	done := make(chan bool)
	sessions := make(chan *base.Session)

	ctx := base.NewContext(config, sessions, done)
	terminator := terminator.NewTerminator(ctx)
	notifier := notifier.NewNotifier(ctx)

	handleSignals(ctx, notifier)

	// Run managers asynchronously and wait for all of them to end
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		terminator.Run()
	}()
	go func() {
		defer wg.Done()
		notifier.Run()
	}()
	wg.Wait()
}

// handleSignals handles operating system signals
func handleSignals(ctx *base.Context, n notifier.Notifier) {
	// When interrupt or terminated, terminate managers, close channel and terminate program
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		for sig := range c {
			log.Debugf("Received %v signal\n", sig)
			close(ctx.Sessions)
			ctx.Done <- true
		}
	}()

	// When hangup, reload notifier
	h := make(chan os.Signal, 1)
	signal.Notify(h, syscall.SIGHUP)
	go func() {
		for sig := range h {
			log.Debugf("Received %v signal\n", sig)
			ctx.Config.Reload()
			n.Reload()
		}
	}()
}

// writePid writes current pid into a pid file
func writePid(file string) {
	log.Infof("Creating pid file %s", file)
	pid := strconv.Itoa(os.Getpid())

	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0644)
	base.Panic(err)
	defer f.Close()

	_, err = f.WriteString(pid)
	base.Panic(err)
}

// removePid removes pid file
func removePid(file string) {
	if _, err := os.Stat(file); err == nil {
		log.Infof("Removing pid file %s", file)
		err := os.Remove(file)
		base.Panic(err)
	}
}

func showVersion() {
	if GitCommit != "" {
		AppVersion = fmt.Sprintf("%s-%s", AppVersion, GitCommit)
	}
	fmt.Printf("%s version %s (compiled with %s)\n", AppName, AppVersion, GoVersion)
}
