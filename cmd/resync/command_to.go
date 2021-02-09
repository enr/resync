package main

import (
	"fmt"
	"path"
	"strings"

	"github.com/enr/go-files/files"
	"github.com/urfave/cli/v2"
)

var commandTo = cli.Command{
	Name: "to",
	// Aliases: []string{"a"},
	Usage:       "",
	Description: `Synchronize all registered folders`,
	Action:      doTo,
	Before:      setupLoadingConfig,
	Flags:       commonFlags,
}

func doTo(c *cli.Context) error {
	args := c.Args().Slice()
	if len(args) < 1 {
		cli.ShowAppHelp(c)
		return exitErrorf(argsNumberErrorExitCode, errorMessageArgsNumber)
	}
	var source string
	var target string
	var err error
	failed := false
	exitCode := 0
	additionalMessage := ""
	baseTarget := strings.TrimSpace(args[0])
	noop := c.Bool("noop")

	if baseTarget == "" || !files.IsDir(baseTarget) {
		return exitErrorf(3, "Invalid base target %s", baseTarget)
	}
	for name, folder := range configuration.Folders {
		ui.Titlef(name)
		source, err = sourcePath(folder.LocalPath)
		if err != nil {
			failed = true
			exitCode = 1
			additionalMessage = fmt.Sprintf("Error resolving source path %s: %v", folder.LocalPath, err)
			break
		}
		target = path.Join(baseTarget, folder.ExternalPath)
		if !files.IsDir(target) {
			ui.Warnf("Target directory %s not found: SKIP", target)
			continue
		}
		options, err := rawSyncOptions(source, target)
		success, code := manageOptionsError(err)
		if !success {
			failed = true
			exitCode = code
			break
		}
		options.Noop = noop
		options.Timestamp = timestamp()
		syncError := syncronizeDirectories(options)
		err = writeReport(options, syncError)
		if syncError != nil {
			failed = true
			exitCode = 1
			additionalMessage = fmt.Sprintf("Exit 1, cause: %v", syncError)
			break
		}
		if err != nil {
			failed = true
			exitCode = 5
			additionalMessage = fmt.Sprintf("Error writing report: %v", err)
			break
		}
	}
	if failed {
		m := "Something gone wrong"
		if additionalMessage != "" {
			m = additionalMessage
		}
		return exitErrorf(exitCode, m)
	}
	return nil
}
