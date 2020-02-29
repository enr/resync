package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/enr/resync/lib/core"
	"github.com/urfave/cli/v2"
)

var (
	runningProcess  *os.Process
	versionTemplate = `%s
Revision: %s
Build date: %s
`
	appVersion  = fmt.Sprintf(versionTemplate, core.Version, core.GitCommit, core.BuildTime)
	commonFlags = []cli.Flag{
		&cli.BoolFlag{Name: "noop", Aliases: []string{"n"}, Usage: "noop mode: no changes made"},
		&cli.BoolFlag{Name: "debug", Aliases: []string{"d"}, Usage: "operates in debug mode: lot of output"},
		&cli.BoolFlag{Name: "quiet", Aliases: []string{"q"}, Usage: "operates in quiet mode"},
		&cli.StringFlag{Name: "config", Aliases: []string{"c"}, Usage: "path to configuration file"},
	}
)

func main() {
	setupCloseHandler()
	realMain(os.Args)
}

func realMain(args []string) {
	app := cli.NewApp()
	app.Name = "resync"
	app.Version = appVersion
	app.Usage = "Synchronize directories using platform tools."
	app.Flags = commonFlags
	app.EnableBashCompletion = true
	app.Action = defaultCommand
	app.Commands = Commands

	app.Run(args)
}

func setupCloseHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Printf("\nStopping %v \n", runningProcess)
		if runningProcess != nil {
			runningProcess.Kill()
		}
		os.Exit(0)
	}()
}
