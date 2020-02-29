package main

import (
	"github.com/urfave/cli/v2"
)

// Default command synchronize a single folder
func defaultCommand(c *cli.Context) error {
	setup(c)
	noop := c.Bool("noop")
	args := c.Args().Slice()
	if len(args) < 2 {
		cli.ShowAppHelp(c)
		return exitErrorf(argsNumberErrorExitCode, errorMessageArgsNumber)
	}
	options, err := rawSyncOptions(args[0], args[1])
	success, exitCode := manageOptionsError(err)
	if !success {
		return exitErrorf(exitCode, "Command execution fail (exit code %d)", exitCode)
	}
	options.Noop = noop
	options.Timestamp = timestamp()

	syncError := syncronizeDirectories(options)
	err = writeReport(options, syncError)
	if syncError != nil {
		return exitErrorf(1, "Exit 1, cause: %v", syncError)
	}
	if err != nil {
		return exitErrorf(5, "Error writing report: %v", err)
	}
	return nil
}
