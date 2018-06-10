package main

import (
	"flag"
	"fmt"
	"github.com/jouir/pgterminate/base"
	"github.com/jouir/pgterminate/notifier"
	"github.com/jouir/pgterminate/terminator"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

// AppVersion stores application version at compilation time
var AppVersion string

func main() {
	var err error
	config := base.NewConfig()

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
	flag.StringVar(&config.LogFile, "log-file", "", "Write logs to a file")
	flag.StringVar(&config.PidFile, "pid-file", "", "Write process id into a file")
	flag.Parse()

	if *version {
		if AppVersion == "" {
			AppVersion = "unknown"
		}
		fmt.Println(AppVersion)
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
		log.Fatalln("Parameter -active-timeout or -idle-timeout required")
	}

	if config.PidFile != "" {
		writePid(config.PidFile)
		defer removePid(config.PidFile)
	}

	done := make(chan bool)
	sessions := make(chan base.Session)

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
			log.Printf("Received %v signal\n", sig)
			close(ctx.Sessions)
			ctx.Done <- true
		}
	}()

	// When hangup, reload notifier
	h := make(chan os.Signal, 1)
	signal.Notify(h, syscall.SIGHUP)
	go func() {
		for sig := range h {
			log.Printf("Received %v signal\n", sig)
			ctx.Config.Reload()
			n.Reload()
		}
	}()
}

// writePid writes current pid into a pid file
func writePid(file string) {
	log.Println("Creating pid file", file)
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
		log.Println("Removing pid file", file)
		err := os.Remove(file)
		base.Panic(err)
	}
}
