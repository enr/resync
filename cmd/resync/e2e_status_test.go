package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/enr/resync/lib/core"
	"github.com/enr/runcmd"
)

// Representation of a command execution.
type commandExecution struct {
	command  *runcmd.Command
	success  bool
	exitCode int
	stdout   string
	stderr   string
}

// resolves the path to the built executable used in tests.
func exePath(p string) string {
	adjusted := fmt.Sprintf("../../bin/%s", p)
	a, _ := filepath.Abs(adjusted)
	var ext string
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}
	executablePath := fmt.Sprintf("%s%s", a, ext)
	if _, err := os.Stat(executablePath); os.IsNotExist(err) {
		panic(fmt.Sprintf("no such executable: %s", executablePath))
	}
	return executablePath
}

var executions = []commandExecution{
	{
		command: &runcmd.Command{
			Exe:  exePath("resync"),
			Args: []string{"-c", fmt.Sprintf("%s/resyncrc.yaml", testdataPath)},
		},
		success:  false,
		exitCode: argsNumberErrorExitCode,
		stderr:   errorMessageArgsNumber,
	},
	{
		command: &runcmd.Command{
			Exe:  exePath("resync"),
			Args: []string{"-c", fmt.Sprintf("%s/resyncrc.yaml", testdataPath), testdataPath, ""},
		},
		success:  false,
		exitCode: 3,
	},
	{
		command: &runcmd.Command{
			Exe:  exePath("resync"),
			Args: []string{"", "."},
		},
		success:  false,
		exitCode: 3,
	},
	{
		command: &runcmd.Command{
			Exe:  exePath("resync"),
			Args: []string{"-c", fmt.Sprintf("%s/resyncrc.yaml", testdataPath), "nosuchdir", "."},
		},
		success:  false,
		exitCode: 4,
	},
	{
		command: &runcmd.Command{
			Exe:  exePath("resync"),
			Args: []string{"-c", fmt.Sprintf("%s/resyncrc.yaml", testdataPath), ".", "nosuchdir"},
		},
		success:  false,
		exitCode: 4,
	},
	{
		command: &runcmd.Command{
			Exe:  exePath("resync"),
			Args: []string{"--version"},
		},
		success:  true,
		exitCode: 0,
		stdout:   fmt.Sprintf("resync version %s", core.Version),
	},
}

func TestCommandExecution(t *testing.T) {
	for _, d := range executions {
		command := d.command
		res := command.Run()
		if res.Success() != d.success {
			t.Fatalf("%s: expected success %t but got %t", command, d.success, res.Success())
		}
		expectedCode := d.exitCode
		actualCode := res.ExitStatus()
		if actualCode != expectedCode {
			t.Fatalf("%s: expected exit code %d but got %d", command, expectedCode, actualCode)
		}
		assertStringContains(t, res.Stdout().String(), d.stdout)
		assertStringContains(t, res.Stderr().String(), d.stderr)
	}
}
