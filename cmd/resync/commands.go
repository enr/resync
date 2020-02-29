package main

import (
	"fmt"

	"github.com/enr/clui"
	"github.com/urfave/cli/v2"
)

var (
	ui            *clui.Clui
	configuration Configuration
)

func setup(c *cli.Context) error {
	if ui != nil {
		return nil
	}
	var err error
	verbosityLevel := clui.VerbosityLevelMedium
	if c.Bool("debug") {
		verbosityLevel = clui.VerbosityLevelHigh
	}
	if c.Bool("quiet") {
		verbosityLevel = clui.VerbosityLevelLow
	}
	ui, err = clui.NewClui(func(ui *clui.Clui) {
		ui.VerbosityLevel = verbosityLevel
	})
	return err
}

func setupLoadingConfig(c *cli.Context) error {
	err := setup(c)
	if err != nil {
		return err
	}
	var cfg string
	if c.String("config") != "" {
		cfg = c.String("config")
		// qua stessi controlli che sotto
		fmt.Printf("file config -c %s \n", cfg)
	} else {
		cfg = configurationFilePath()
	}
	configuration, err = loadConfigurationFromFile(cfg)
	return err
}

// Commands other than default action.
var Commands = []*cli.Command{
	//	commandFrom,
	&commandTo,
}

func exitErrorf(exitCode int, template string, args ...interface{}) error {
	ui.Errorf(`Something gone wrong:`)
	return cli.NewExitError(fmt.Sprintf(template, args...), exitCode)
}
